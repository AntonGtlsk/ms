package get

import (
	u "auth-ms/internal/entity/user"
	resp "auth-ms/internal/lib/api/responce"
	"auth-ms/internal/service/auth"
	"logging"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type Request struct {
	UserId string `json:"userid" validate:"required"`
}

type Response struct {
	resp.Response `json:"response"`
	Subscription  u.Subscription `json:"subscription"`
}

type SubGetter interface {
	GetSubscription(userid string) (u.Subscription, error)
}

func New(logging logging.LoggerInterface, subGetter SubGetter, jwtParser *auth.JwtParser) http.HandlerFunc {
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

		var req Request

		req.UserId = r.URL.Query().Get("userid")

		if err := validator.New().Struct(req); err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("Failed to validate data")

			validateErr := err.(validator.ValidationErrors)

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		if req.UserId != uid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		subscription, err := subGetter.GetSubscription(req.UserId)

		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("Failed to getting subscription %s: ", req.UserId)

			render.JSON(w, r, resp.Error("Failed to getting subscription"))

			return
		}

		render.JSON(w, r, Response{Subscription: subscription, Response: resp.OK()})

		logging.Infof("Successful getting of the subscription")

	}
}
