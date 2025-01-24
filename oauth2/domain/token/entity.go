package token

import (
	"time"

	"github.com/google/uuid"
)

func NewToken(
	accessToken string,
	clientID uuid.UUID,
	userID uuid.UUID,
	scope string,
	tokenExpiresAt time.Time,
	refreshExpiresAt time.Time) *Token {
	return &Token{
		AccessToken: accessToken,
		ClientID:    clientID,
		UserID:      userID,
		Scope:       scope,
		ExpiresAt:   tokenExpiresAt,
	}
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
