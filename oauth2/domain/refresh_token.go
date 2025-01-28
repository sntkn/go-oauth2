package domain

import (
	"time"

	"github.com/sntkn/go-oauth2/oauth2/pkg/str"
)

const (
	randomStringLen = 32
	day             = 24 * time.Hour
)

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

type StoreNewRefreshTokenParams struct {
	AccessToken   string
	AdditionalDay int
	Repo          RefreshTokenRepository
}

func StoreNewRefreshToken(p StoreNewRefreshTokenParams) (RefreshToken, error) {
	rtoken := NewRefreshToken(RefreshTokenParams{
		AccessToken: p.AccessToken,
	})
	if err := rtoken.SetNewRefreshToken(); err != nil {
		return nil, err
	}
	rtoken.SetNewExpiry(p.AdditionalDay)

	if err := p.Repo.StoreRefreshToken(rtoken); err != nil {
		return nil, err
	}
	return rtoken, nil
}

//go:generate go run github.com/matryer/moq -out refresh_token_mock.go . RefreshToken
type RefreshToken interface {
	IsNotFound() bool
	GetRefreshToken() string
	GetAccessToken() string
	GetExpiresAt() time.Time
	Expiry() int64
	SetNewRefreshToken() error
	SetNewExpiry(additionalDays int)
}

//go:generate go run github.com/matryer/moq -out refresh_token_repository_mock.go . RefreshTokenRepository
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

func (t *refreshToken) IsNotFound() bool {
	return t.RefreshToken == ""
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

func (t *refreshToken) SetNewRefreshToken() error {
	randomString, err := str.GenerateRandomString(randomStringLen)
	if err != nil {
		return err
	}

	t.RefreshToken = randomString

	return nil
}

func (t *refreshToken) SetNewExpiry(additionalDays int) {
	t.ExpiresAt = time.Now().Add(time.Duration(additionalDays) * day)
}
