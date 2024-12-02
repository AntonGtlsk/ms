package update

import (
	u "auth-ms/internal/entity/user"
	resp "auth-ms/internal/lib/api/responce"
	cv "auth-ms/internal/lib/validator"
	"auth-ms/internal/service/auth"
	"fmt"
	"logging"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type Request struct {
	UserId       string         `json:"userid" validate:"required"`
	Subscription u.Subscription `json:"subscription" validate:"required,isSub"`
}

type Response struct {
	resp.Response `json:"response"`
	Msg           string `json:"message"`
}

type SubUpdater interface {
	UpdateSubscription(userid string, subscription u.Subscription) (int64, error)
}

func New(logging logging.LoggerInterface, subUpdater SubUpdater, jwtParser auth.JwtParser) http.HandlerFunc {
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

		err = render.DecodeJSON(r.Body, &req)

		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("Failed to decode request")

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		if req.UserId != uid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		v := validator.New()

		v.RegisterValidation("isSub", cv.IsSubValidator)

		if err := v.Struct(req); err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("Failed to validate data")

			validateErr := err.(validator.ValidationErrors)

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		rowsAffected, err := subUpdater.UpdateSubscription(req.UserId, req.Subscription)

		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("Failed to update subscription %s: %v", req.UserId, req.Subscription)

			render.JSON(w, r, resp.Error(fmt.Sprintf("failed to update subscription: %v", err)))

			return
		}

		if rowsAffected == 0 {
			responseOK(w, r, "everything is updated")
		} else {
			responseOK(w, r, fmt.Sprintf("success", rowsAffected))
		}

	}
}

func responseOK(w http.ResponseWriter, r *http.Request, msg string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Msg:      msg,
	})
}
