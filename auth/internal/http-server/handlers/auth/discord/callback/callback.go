package callback

import (
	ea "auth-ms/internal/entity/auth"
	u "auth-ms/internal/entity/user"
	"auth-ms/internal/service/auth"
	"fmt"
	"logging"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-chi/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

const (
	discordClientID     = ""
	discordClientSecret = ""
	discordRedirectURL  = "http://localhost:80/login"
)

type DiscordResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type Response struct {
	AccessToken string `json:"access_token"`
	User        u.User `json:"user"`
}

type UserProvider interface {
	IsUserExist(userId string) (bool, error)
	SaveUser(usr u.SaveUser) (u.User, bool, error)
	UpdateUser(user u.UpdateUser) (u.User, int64, error)
	SetRefreshSession(userId string, refreshToken ea.RefreshSession) error
}

func New(logging logging.LoggerInterface, userProvider UserProvider, jwtParser *auth.JwtParser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

		data := url.Values{}
		data.Set("grant_type", "authorization_code")
		data.Set("code", code)
		data.Set("redirect_uri", discordRedirectURL)

		client := &http.Client{}
		req, err := http.NewRequest("POST", "https://discord.com/api/oauth2/token", strings.NewReader(data.Encode()))
		if err != nil {
			fmt.Println("Ошибка создания запроса:", err)
			return
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.SetBasicAuth(discordClientID, discordClientSecret)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Ошибка выполнения запроса:", err)
			return
		}
		defer resp.Body.Close()

		var discordResponse DiscordResponse
		err = render.DecodeJSON(resp.Body, &discordResponse)

		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}

		discord, err := discordgo.New("Bearer " + discordResponse.AccessToken)
		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("Failed to create Discord session")
			http.Error(w, "Failed to create Discord session", http.StatusInternalServerError)
			return
		}

		user, err := discord.User("@me")
		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("Failed to fetch user information")
			http.Error(w, "Failed to fetch user information", http.StatusInternalServerError)
			return
		}

		exist, err := userProvider.IsUserExist(user.ID)

		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("Failed to fetch user information")
			http.Error(w, "Failed to fetch user information", http.StatusInternalServerError)
			return
		}

		var respUser u.User

		if exist {

			respUser, _, err = userProvider.UpdateUser(u.UpdateUser{
				UserId:    user.ID,
				Name:      user.Username,
				AvatarURL: user.AvatarURL(""),
			})

			if err != nil {
				logging.WithFields(logrus.Fields{
					"error": err,
				}).Errorf("Failed to update user")
			}

		} else {
			respUser, _, err = userProvider.SaveUser(u.SaveUser{
				UserId:    user.ID,
				Name:      user.Username,
				AvatarURL: user.AvatarURL(""),
				Sub:       u.None,
			})

			if err != nil {
				logging.WithFields(logrus.Fields{
					"error": err,
				}).Errorf("Failed to save user")
			}
		}

		tokenString, err := jwtParser.GenerateJWTToken(float64(respUser.Id), user.ID, user.Username, user.AvatarURL(""), respUser.Sub)
		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("Failed to generate jwt token")
			http.Error(w, "Failed to generate jwt token", http.StatusInternalServerError)
			return
		}

		refreshToken, claims, err := jwtParser.GenerateRefreshToken(float64(respUser.Id), user.ID, user.Username, user.AvatarURL(""))
		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("Failed to generate refresh token")
			http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
			return
		}

		err = userProvider.SetRefreshSession(user.ID, ea.RefreshSession{
			RefreshToken:        refreshToken,
			DiscordRefreshToken: discordResponse.RefreshToken,
			ExpiresIn:           claims["exp"].(int64),
			CreatedAt:           claims["createdAt"].(int64),
		},
		)
		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("Failed to save refresh session")
			http.Error(w, "Failed to save refresh session", http.StatusInternalServerError)
			return
		}

		cookie := http.Cookie{
			Name:     "Refresh-Token",
			Path:     "",
			Domain:   "localhost",
			Value:    refreshToken,
			Expires:  time.Now().Add(time.Hour * 24 * 30),
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)

		var response Response

		response.AccessToken = tokenString
		response.User = respUser

		render.JSON(w, r, response)
	}
}
