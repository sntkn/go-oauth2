package validation

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func RequiredWithFieldValue(fl validator.FieldLevel) bool {
	params := strings.Split(fl.Param(), " ")
	if len(params) != 2 {
		return false
	}
	targetFieldValue := fl.Parent().FieldByName(params[0])
	if targetFieldValue.IsValid() && targetFieldValue.String() == params[1] {
		return fl.Field().String() != ""
	}
	return true
}
