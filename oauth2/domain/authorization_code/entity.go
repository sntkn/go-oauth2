package authorization_code

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/pkg/str"
)

func NewAuthorizationCode(code string, clientID, userID uuid.UUID, scope, redirectURI string, expiresAt, createdAt, UpdatedAt time.Time) *AuthorizationCode {
	return &AuthorizationCode{
		Code:        code,
		ClientID:    clientID,
		UserID:      userID,
		Scope:       scope,
		RedirectURI: redirectURI,
		ExpiresAt:   expiresAt,
		CreatedAt:   createdAt,
		UpdatedAt:   UpdatedAt,
	}
}

type AuthorizationCode struct {
	Code     string
	ClientID uuid.UUID
	UserID   uuid.UUID
	Scope    string
	// Scopes      []string
	RedirectURI string
	ExpiresAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (a *AuthorizationCode) GenerateRedirectURIWithCode() string {
	return fmt.Sprintf("%s?code=%s", a.RedirectURI, a.Code)
}

func GenerateCode() (string, error) {
	randomStringLen := 32
	return str.GenerateRandomString(randomStringLen)
}
