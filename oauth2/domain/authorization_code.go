package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/pkg/str"
)

type StoreAuthorizationCodeParams struct {
	Code        string
	ClientID    uuid.UUID
	UserID      uuid.UUID
	Scope       string
	RedirectURI string
	ExpiresAt   time.Time
}

type AuthorizationCodeParams struct {
	Code        string
	ClientID    string
	UserID      string
	Scope       string
	RedirectURI string
	ExpiresAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewAuthorizationCode(p AuthorizationCodeParams) (AuthorizationCode, error) {

	clientID, err := uuid.Parse(p.ClientID)
	if err != nil {
		return nil, err
	}
	userID, err := uuid.Parse(p.UserID)
	if err != nil {
		return nil, err
	}

	return &authorizationCode{
		Code:        p.Code,
		ClientID:    clientID,
		UserID:      userID,
		Scope:       p.Scope,
		RedirectURI: p.RedirectURI,
		ExpiresAt:   p.ExpiresAt,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}, nil
}

//go:generate go run github.com/matryer/moq -out authorization_code_mock.go . AuthorizationCode
type AuthorizationCode interface {
	IsNotFound() bool
	GetCode() string
	GetClientID() uuid.UUID
	GetUserID() uuid.UUID
	GetScope() string
	GetRedirectURI() string
	GetExpiresAt() time.Time
	GenerateRedirectURIWithCode() string
	IsExpired(t time.Time) bool
}

//go:generate go run github.com/matryer/moq -out authorization_code_repository_mock.go . AuthorizationCodeRepository
type AuthorizationCodeRepository interface {
	FindAuthorizationCode(string) (AuthorizationCode, error)
	StoreAuthorizationCode(StoreAuthorizationCodeParams) (string, error)
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

func (a *authorizationCode) IsExpired(t time.Time) bool {
	return t.After(a.ExpiresAt)
}

func GenerateCode() (string, error) {
	randomStringLen := 32
	return str.GenerateRandomString(randomStringLen)
}
