package authorization

import (
	"time"

	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/infrastructure/model"
)

type IAuthorizationRepository interface {
	FindClientByClientID(uuid.UUID) (*model.Client, error)
	FindAuthorizationCode(string) (*model.AuthorizationCode, error)
	StoreAuthorizationCode(*model.AuthorizationCode) error
	FindValidAuthorizationCode(string, time.Time) (*model.AuthorizationCode, error)
	StoreToken(*model.Token) error
	StoreRefreshToken(t *model.RefreshToken) error
	RevokeCode(code string) error
}
