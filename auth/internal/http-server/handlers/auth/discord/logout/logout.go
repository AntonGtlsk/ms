package logout

import (
	resp "auth-ms/internal/lib/api/responce"
	"logging"
	"net/http"

	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

type RefreshSessionDeleter interface {
	DeleteRefreshSession(token string) (int64, error)
}

type Response struct {
	resp.Response `json:"response"`
	Msg           string `json:"message"`
}

func New(logging logging.LoggerInterface, deleter RefreshSessionDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Refresh-Token")

		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("Failed to get cookie")
			http.Error(w, "Failed to get cookie", http.StatusInternalServerError)
			return
		}

		refreshToken := cookie.Value

		_, err = deleter.DeleteRefreshSession(refreshToken)

		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("Failed to delete session")
			http.Error(w, "Failed to logout", http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, Response{Response: resp.OK(), Msg: "Successful logout"})

	}
}
