package refresh

import (
	ea "auth-ms/internal/entity/auth"
	u "auth-ms/internal/entity/user"
	"auth-ms/internal/service/auth"
	"logging"
	"time"

	"net/http"

	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

type RefreshSessionAndUserManager interface {
	SetRefreshSession(userId string, refreshSession ea.RefreshSession) error
	GetRefreshSession(token string) (ea.RefreshSession, error)
	DeleteRefreshSession(token string) (int64, error)
	GetSubscription(userid string) (u.Subscription, error)
}

type Response struct {
	AccessToken string `json:"access_token"`
}

func New(logging logging.LoggerInterface, refreshSessionAndUserManager RefreshSessionAndUserManager, jwtParser *auth.JwtParser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Refresh-Token")
		r.Cookies()
		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("Failed to get cookie")
			http.Error(w, "Failed to get cookie", http.StatusInternalServerError)
			return
		}

		refreshToken := cookie.Value

		id, userid, username, avatarURL, _, _, err := jwtParser.ParseRefreshToken(refreshToken)

		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("Failed to parse refresh token")
			http.Error(w, "Failed to parse refresh token", http.StatusInternalServerError)
			return
		}

		subscription, err := refreshSessionAndUserManager.GetSubscription(userid)

		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("Failed to parse get subscription")
			http.Error(w, "Failed to parse get subscription", http.StatusInternalServerError)
			return
		}

		newAccessToken, err := jwtParser.GenerateJWTToken(float64(id), userid, username, avatarURL, subscription)
		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("Failed to generate jwt token")
			http.Error(w, "Failed to generate jwt token", http.StatusInternalServerError)
			return
		}
		newRefreshToken, claims, err := jwtParser.GenerateRefreshToken(float64(id), userid, username, avatarURL)

		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("Failed to generate refresh token")
			http.Error(w, "Failed to sagenerateve refresh token", http.StatusInternalServerError)
			return
		}

		_, err = refreshSessionAndUserManager.GetRefreshSession(refreshToken)

		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("Failed to get refresh session")
			http.Error(w, "Failed to get refresh session", http.StatusInternalServerError)
			return
		}

		_, err = refreshSessionAndUserManager.DeleteRefreshSession(refreshToken)

		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("Failed to delete old refresh session")
			http.Error(w, "Failed to delete old refresh session", http.StatusInternalServerError)
			return
		}

		err = refreshSessionAndUserManager.SetRefreshSession(userid, ea.RefreshSession{
			RefreshToken: newRefreshToken,
			ExpiresIn:    claims["exp"].(int64),
			CreatedAt:    claims["createdAt"].(int64),
		})

		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("Error calling refreshSessionAndUserManager.SetRefreshSession")
		}

		newCookie := http.Cookie{
			Name:     "Refresh-Token",
			Path:     "",
			Domain:   "localhost",
			Value:    newRefreshToken,
			Expires:  time.Now().Add(time.Hour * 24 * 30),
			HttpOnly: true,
		}
		http.SetCookie(w, &newCookie)

		var resp Response

		resp.AccessToken = newAccessToken
		render.JSON(w, r, resp)
	}
}
