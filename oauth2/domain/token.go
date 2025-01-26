package domain

import (
	"time"

	"github.com/google/uuid"
)

type TokenParams struct {
	AccessToken string
	ClientID    uuid.UUID
	UserID      uuid.UUID
	Scope       string
	ExpiresAt   time.Time
}

func NewToken(p TokenParams) Token {
	return &token{
		AccessToken: p.AccessToken,
		ClientID:    p.ClientID,
		UserID:      p.UserID,
		Scope:       p.Scope,
		ExpiresAt:   p.ExpiresAt,
	}
}

type Token interface {
	GetAccessToken() string
	GetClientID() uuid.UUID
	GetUserID() uuid.UUID
	GetScope() string
	GetExpiresAt() time.Time
	Expiry() int64
}

type TokenRepository interface {
	StoreToken(Token) error
	FindToken(accessToken string) (Token, error)
	RevokeToken(accessToken string) error
}

type token struct {
	AccessToken string
	ClientID    uuid.UUID
	UserID      uuid.UUID
	Scope       string
	ExpiresAt   time.Time
}

func (t *token) GetAccessToken() string {
	return t.AccessToken
}

func (t *token) GetClientID() uuid.UUID {
	return t.ClientID
}

func (t *token) GetUserID() uuid.UUID {
	return t.UserID
}

func (t *token) GetScope() string {
	return t.Scope
}

func (t *token) GetExpiresAt() time.Time {
	return t.ExpiresAt
}

func (t *token) SetExpiresAt(tim time.Time) {
	t.ExpiresAt = tim
}

func (t *token) Expiry() int64 {
	return t.ExpiresAt.Unix()
}
