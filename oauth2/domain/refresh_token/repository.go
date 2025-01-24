package refresh_token

import (
	"time"
)

type IAuthorizationRepository interface {
	StoreRefreshToken(t *RefreshToken) error
	FindValidRefreshToken(refreshToken string, expiresAt time.Time) (*RefreshToken, error)
	RevokeRefreshToken(refreshToken string) error
}
