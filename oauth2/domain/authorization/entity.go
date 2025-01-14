package authorization

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) IsNotFound() bool {
	return u.ID == uuid.Nil
}

func (u *User) IsPasswordMatch(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

func NewUser(id uuid.UUID, name, email, password string, createdAt, updatedAt time.Time) *User {
	return &User{
		ID:        id,
		Name:      name,
		Email:     email,
		Password:  password,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

type Client struct {
	ID           uuid.UUID
	Name         string
	RedirectURIs string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (c *Client) IsNotFound() bool {
	return c.ID == uuid.Nil
}

func (c *Client) IsRedirectURIMatch(redirectURI string) bool {
	// TODO: 複数のリダイレクトURIを持つ場合の対応
	return c.RedirectURIs == redirectURI
}

func NewClient(id uuid.UUID, name, redirectURIs string, createdAt, updatedAt time.Time) *Client {
	return &Client{
		ID:           id,
		Name:         name,
		RedirectURIs: redirectURIs,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}

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
