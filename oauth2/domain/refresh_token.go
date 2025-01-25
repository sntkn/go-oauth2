package domain

import "time"

type RefreshTokenRepository interface {
	StoreRefreshToken(t *RefreshToken) error
	FindValidRefreshToken(refreshToken string, expiresAt time.Time) (*RefreshToken, error)
	RevokeRefreshToken(refreshToken string) error
}

type RefreshToken struct {
	RefreshToken string
	AccessToken  string
	ExpiresAt    time.Time
}
