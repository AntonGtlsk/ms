package validator

import (
	"github.com/go-playground/validator/v10"
)

func IsSubValidator(fl validator.FieldLevel) bool {
	sub := fl.Field().String()

	switch sub {
	case "none":
		return true
	case "sub":
		return true
	default:
		return false
	}
}
