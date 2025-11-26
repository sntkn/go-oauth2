package domain

import (
	"context"
	"time"

	"github.com/sntkn/go-oauth2/oauth2/pkg/str"
)

const (
	randomStringLen = 32
	day             = 24 * time.Hour
)

type RefreshTokenParams struct {
	RefreshToken RefreshTokenString
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

//go:generate go run github.com/matryer/moq -out refresh_token_mock.go . RefreshToken
type RefreshToken interface {
	IsNotFound() bool
	GetRefreshToken() string
	GetAccessToken() string
	GetExpiresAt() time.Time
	Expiry() int64
	SetNewExpiry(additionalDays int)
	IsExpired(now time.Time) bool
}

//go:generate go run github.com/matryer/moq -out refresh_token_repository_mock.go . RefreshTokenRepository
type RefreshTokenRepository interface {
	StoreRefreshToken(ctx context.Context, t RefreshToken) error
	FindRefreshToken(ctx context.Context, refreshToken string) (RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, refreshToken string) error
}

type refreshToken struct {
	RefreshToken RefreshTokenString
	AccessToken  string
	ExpiresAt    time.Time
}

func (t *refreshToken) IsNotFound() bool {
	return t.RefreshToken == ""
}

func (t *refreshToken) GetRefreshToken() string {
	return t.RefreshToken.String()
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

func (t *refreshToken) IsExpired(now time.Time) bool {
	return t.ExpiresAt.After(now)
}

func (t *refreshToken) SetNewExpiry(additionalDays int) {
	t.ExpiresAt = time.Now().Add(time.Duration(additionalDays) * day)
}

type RefreshTokenString string

func (s RefreshTokenString) String() string {
	return string(s)
}

func (_ RefreshTokenString) Generate() (RefreshTokenString, error) {
	randomString, err := str.GenerateRandomString(randomStringLen)
	if err != nil {
		return "", err
	}

	return RefreshTokenString(randomString), nil
}
