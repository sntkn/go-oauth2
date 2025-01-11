package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Client struct {
	ID           uuid.UUID `db:"id"`
	Name         string    `db:"name"`
	RedirectURIs string    `db:"redirect_uris"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type Code struct {
	Code        string    `db:"code"`
	ClientID    uuid.UUID `db:"client_id"`
	UserID      uuid.UUID `db:"user_id"`
	Scope       string    `db:"scope"`
	RedirectURI string    `db:"redirect_uri"`
	ExpiresAt   time.Time `db:"expires_at"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type Token struct {
	AccessToken string    `db:"access_token"`
	ClientID    uuid.UUID `db:"client_id"`
	UserID      uuid.UUID `db:"user_id"`
	Scope       string    `db:"scope"`
	ExpiresAt   time.Time `db:"expires_at"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type RefreshToken struct {
	RefreshToken string    `db:"refresh_token"`
	AccessToken  string    `db:"access_token"`
	ExpiresAt    time.Time `db:"expires_at"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
