package domain

import "time"

type RefreshTokenParams struct {
	RefreshToken string
	AccessToken  string
	ExpiresAt    time.Time
}

func NewRefreshToken(p RefreshTokenParams) RefreshToken {
	return &refreshToken{
		RefreshToken: p.RefreshToken,
		AccessToken:  p.AccessToken,
		ExpiresAt:    p.ExpiresAt,
	}
}

type RefreshToken interface {
	GetRefreshToken() string
	GetAccessToken() string
	GetExpiresAt() time.Time
	Expiry() int64
}

type RefreshTokenRepository interface {
	StoreRefreshToken(t RefreshToken) error
	FindValidRefreshToken(refreshToken string, expiresAt time.Time) (RefreshToken, error)
	RevokeRefreshToken(refreshToken string) error
}

type refreshToken struct {
	RefreshToken string
	AccessToken  string
	ExpiresAt    time.Time
}

func (t *refreshToken) GetRefreshToken() string {
	return t.RefreshToken
}

func (t *refreshToken) GetAccessToken() string {
	return t.AccessToken
}

func (t *refreshToken) GetExpiresAt() time.Time {
	return t.ExpiresAt
}

func (t *refreshToken) Expiry() int64 {
	return t.ExpiresAt.Unix()
}
