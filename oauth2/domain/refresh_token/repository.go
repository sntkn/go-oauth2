package refresh_token

import (
	"time"
)

type RefreshTokenRepository interface {
	StoreRefreshToken(t *RefreshToken) error
	FindValidRefreshToken(refreshToken string, expiresAt time.Time) (*RefreshToken, error)
	RevokeRefreshToken(refreshToken string) error
}
