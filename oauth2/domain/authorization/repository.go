package authorization

import (
	"time"

	"github.com/google/uuid"
)

type IAuthorizationRepository interface {
	FindClientByClientID(uuid.UUID) (*Client, error)
	FindAuthorizationCode(string) (*AuthorizationCode, error)
	StoreAuthorizationCode(*AuthorizationCode) error
	FindValidAuthorizationCode(string, time.Time) (*AuthorizationCode, error)
	StoreToken(*Token) error
	StoreRefreshToken(t *RefreshToken) error
	RevokeCode(code string) error
	FindValidRefreshToken(refreshToken string, expiresAt time.Time) (*RefreshToken, error)
	FindToken(accessToken string) (*Token, error)
	RevokeToken(accessToken string) error
	RevokeRefreshToken(refreshToken string) error
}
