package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/pkg/str"
)

type AuthorizationCodeParams struct {
	Code        string
	ClientID    uuid.UUID
	UserID      uuid.UUID
	Scope       string
	RedirectURI string
	ExpiresAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewAuthorizationCode(p AuthorizationCodeParams) *authorizationCode {
	return &authorizationCode{
		Code:        p.Code,
		ClientID:    p.ClientID,
		UserID:      p.UserID,
		Scope:       p.Scope,
		RedirectURI: p.RedirectURI,
		ExpiresAt:   p.ExpiresAt,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

type AuthorizationCode interface {
	IsNotFound() bool
	GetCode() string
	GetClientID() uuid.UUID
	GetUserID() uuid.UUID
	GetScope() string
	GetRedirectURI() string
	GetExpiresAt() time.Time
	GenerateRedirectURIWithCode() string
}

type AuthorizationCodeRepository interface {
	FindAuthorizationCode(string) (AuthorizationCode, error)
	StoreAuthorizationCode(AuthorizationCode) error
	FindValidAuthorizationCode(string, time.Time) (AuthorizationCode, error)
	RevokeCode(code string) error
}

type authorizationCode struct {
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

func (a *authorizationCode) IsNotFound() bool {
	return a.Code == ""
}

func (a *authorizationCode) GetCode() string {
	return a.Code
}

func (a *authorizationCode) GetClientID() uuid.UUID {
	return a.ClientID
}

func (a *authorizationCode) GetUserID() uuid.UUID {
	return a.UserID
}

func (a *authorizationCode) GetScope() string {
	return a.Scope
}

func (a *authorizationCode) GetRedirectURI() string {
	return a.RedirectURI
}

func (a *authorizationCode) GetExpiresAt() time.Time {
	return a.ExpiresAt
}

func (a *authorizationCode) GenerateRedirectURIWithCode() string {
	return fmt.Sprintf("%s?code=%s", a.RedirectURI, a.Code)
}

func GenerateCode() (string, error) {
	randomStringLen := 32
	return str.GenerateRandomString(randomStringLen)
}
