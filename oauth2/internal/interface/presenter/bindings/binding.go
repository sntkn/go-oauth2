package bindings

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/sntkn/go-oauth2/oauth2/internal/interface/presenter/validation"
)

// Setup カスタムバリデータを登録
func Setup() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("required_with_field_value", validation.RequiredWithFieldValue); err != nil {
			return err
		}
	}

	return nil
}
