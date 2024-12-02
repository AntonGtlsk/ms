package get_user_guilds

import (
	discordEntity "auth-ms/internal/entity/discord"
	resp "auth-ms/internal/lib/api/responce"
	"auth-ms/internal/service/auth"
	"logging"

	//sa "auth-ms/internal/service/auth"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

const (
	discordClientID     = ""
	discordClientSecret = ""
)

type Request struct {
	UserId string `json:"userid"`
}

type DiscordResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type DiscordRefreshTokenManager interface {
	GetDiscordRefreshToken(userid string) (string, error)
	UpdateDiscordRefreshToken(oldToken, newToken string) (int64, error)
}

func NewOptions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.WriteHeader(http.StatusOK)
		return
	}
}

func New(logging logging.LoggerInterface, discordRefreshTokenManager DiscordRefreshTokenManager, dg *discordgo.Session, jwtParser *auth.JwtParser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		_, uid, _, _, err := jwtParser.ParseToken(tokenString)

		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("Unauthorized")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return

		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		var request Request

		request.UserId = r.URL.Query().Get("userid")

		if request.UserId != uid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if err := validator.New().Struct(request); err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("Failed to validate data")

			validateErr := err.(validator.ValidationErrors)

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		discordRefreshToken, err := discordRefreshTokenManager.GetDiscordRefreshToken(request.UserId)

		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("Failed to get guilds for userid: %s", request.UserId)

			render.JSON(w, r, resp.Error(fmt.Sprintf("failed to get user`s guilds: %v", err)))

			return
		}

		data := url.Values{}
		data.Set("grant_type", "refresh_token")
		data.Set("refresh_token", discordRefreshToken)

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

		_, err = discordRefreshTokenManager.UpdateDiscordRefreshToken(discordRefreshToken, discordResponse.RefreshToken)

		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("Failed to update discord refresh token for user %s: ", request.UserId)

			render.JSON(w, r, fmt.Sprintf("failed to update discord refresh token: %v", err))

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

		userGuilds, err := discord.UserGuilds(100, "", "")
		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("Failed to fetch userGuild information")
			http.Error(w, "Failed to fetch user information", http.StatusInternalServerError)
			return
		}

		botsGuilds, err := dg.UserGuilds(100, "", "")

		var guilds []*discordgo.UserGuild

		for _, botsGuild := range botsGuilds {
			for _, userGuild := range userGuilds {
				if botsGuild.ID == userGuild.ID {
					guilds = append(guilds, userGuild)
				}
			}
		}

		var resultGuilds []discordEntity.Guild

		for _, guild := range guilds {
			members, err := dg.GuildMembers(guild.ID, "", 1000)
			if err != nil {
				logging.WithFields(logrus.Fields{
					"error": err,
				}).Errorf("Error getting guild members with guild id: %s", guild.ID)
			}

			for _, member := range members {
				for _, roleID := range member.Roles {

					role, err := dg.State.Role(guild.ID, roleID)

					if err != nil {
						logging.WithFields(logrus.Fields{
							"error": err,
						}).Errorf("Error getting role with guild id: %s and role id: %s", guild.ID, roleID)
					}

					if err == nil && role.Permissions&discordgo.PermissionAdministrator != 0 {
						if member.User.ID == request.UserId {
							resultGuilds = append(resultGuilds, discordEntity.Guild{
								ID:   guild.ID,
								Name: guild.Name,
							})
						}
					}
				}

			}
		}

		dg.Close()

		render.JSON(w, r, resultGuilds)
	}
}
