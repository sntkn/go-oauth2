package domain

import (
	"time"

	"github.com/google/uuid"
)

type TokenRepository interface {
	StoreToken(*Token) error
	FindToken(accessToken string) (*Token, error)
	RevokeToken(accessToken string) error
}

type Token struct {
	AccessToken string
	ClientID    uuid.UUID
	UserID      uuid.UUID
	Scope       string
	ExpiresAt   time.Time
}

func (t *Token) Expiry() int64 {
	return t.ExpiresAt.Unix()
}
