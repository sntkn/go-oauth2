package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/pkg/str"
)

type uuidLike interface {
	~string | uuid.UUID
}

type StoreAuthorizationCodeParams struct {
	Code        string
	ClientID    uuid.UUID
	UserID      uuid.UUID
	Scope       string
	RedirectURI string
	ExpiresAt   time.Time
}

type AuthorizationCodeParams[T uuidLike] struct {
	Code        string
	ClientID    T
	UserID      T
	Scope       string
	RedirectURI string
	ExpiresAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewAuthorizationCode[T uuidLike](p AuthorizationCodeParams[T]) (AuthorizationCode, error) {
	clientID, err := castUUID(p.ClientID)
	if err != nil {
		return nil, err
	}
	userID, err := castUUID(p.UserID)
	if err != nil {
		return nil, err
	}

	return &authorizationCode{
		code:        p.Code,
		clientID:    clientID,
		userID:      userID,
		scope:       p.Scope,
		redirectURI: p.RedirectURI,
		expiresAt:   p.ExpiresAt,
		createdAt:   p.CreatedAt,
		updatedAt:   p.UpdatedAt,
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
	FindAuthorizationCode(ctx context.Context, code string) (AuthorizationCode, error)
	StoreAuthorizationCode(ctx context.Context, p StoreAuthorizationCodeParams) (string, error)
	FindValidAuthorizationCode(ctx context.Context, code string, expiresAt time.Time) (AuthorizationCode, error)
	RevokeCode(ctx context.Context, code string) error
}

type authorizationCode struct {
	code     string
	clientID uuid.UUID
	userID   uuid.UUID
	scope    string
	// Scopes      []string
	redirectURI string
	expiresAt   time.Time
	createdAt   time.Time
	updatedAt   time.Time
}

func (a *authorizationCode) IsNotFound() bool {
	return a.code == ""
}

func (a *authorizationCode) GetCode() string {
	return a.code
}

func (a *authorizationCode) GetClientID() uuid.UUID {
	return a.clientID
}

func (a *authorizationCode) GetUserID() uuid.UUID {
	return a.userID
}

func (a *authorizationCode) GetScope() string {
	return a.scope
}

func (a *authorizationCode) GetRedirectURI() string {
	return a.redirectURI
}

func (a *authorizationCode) GetExpiresAt() time.Time {
	return a.expiresAt
}

func (a *authorizationCode) GenerateRedirectURIWithCode() string {
	return fmt.Sprintf("%s?code=%s", a.redirectURI, a.code)
}

func (a *authorizationCode) IsExpired(t time.Time) bool {
	return t.After(a.expiresAt)
}

func GenerateCode() (string, error) {
	randomStringLen := 32
	return str.GenerateRandomString(randomStringLen)
}

func castUUID[T uuidLike](v T) (uuid.UUID, error) {
	switch val := any(v).(type) {
	case string:
		return uuid.Parse(val)
	case uuid.UUID:
		return val, nil
	default:
		return uuid.Nil, fmt.Errorf("unsupported uuid type %T", v)
	}
}
