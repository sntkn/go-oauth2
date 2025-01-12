package authorization

import "github.com/sntkn/go-oauth2/oauth2/infrastructure/model"

type IAuthorizationRepository interface {
	FindAuthorizationCode(code string) (*model.Code, error)
	StoreAuthorizationCode(code *model.Code) error
}
