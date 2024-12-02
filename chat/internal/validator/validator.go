package validator

import (
	"strconv"

	"github.com/go-playground/validator/v10"
)

func MaxLengthValidator(fl validator.FieldLevel) bool {
	maxLen, err := strconv.Atoi(fl.Param())

	if err != nil {
		return false
	}

	return len(fl.Field().String()) <= maxLen
}
