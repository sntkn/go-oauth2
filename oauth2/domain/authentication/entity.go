package authentication

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(id uuid.UUID, name, email string, createdAt, updatedAt time.Time) *User {
	return &User{
		ID:        id,
		Name:      name,
		Email:     email,
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

func NewClient(id uuid.UUID, name, redirectURIs string, createdAt, updatedAt time.Time) *Client {
	return &Client{
		ID:           id,
		Name:         name,
		RedirectURIs: redirectURIs,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}
