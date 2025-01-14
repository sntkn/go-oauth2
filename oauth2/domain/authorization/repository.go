package authorization

import (
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/infrastructure/model"
)

type IAuthorizationRepository interface {
	FindClientByClientID(clientID uuid.UUID) (*model.Client, error)
	FindAuthorizationCode(code string) (*model.AuthorizationCode, error)
	StoreAuthorizationCode(code *model.AuthorizationCode) error
}
