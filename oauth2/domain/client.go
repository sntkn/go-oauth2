package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ClientParams struct {
	ID           uuid.UUID
	Name         string
	RedirectURIs []string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewClient(p ClientParams) Client {
	return &client{
		ID:           p.ID,
		Name:         p.Name,
		RedirectURIs: p.RedirectURIs,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}

//go:generate go run github.com/matryer/moq -out client_mock.go . Client
type Client interface {
	IsNotFound() bool
	IsRedirectURIMatch(redirectURI string) bool
}

//go:generate go run github.com/matryer/moq -out client_repository_mock.go . ClientRepository
type ClientRepository interface {
	FindClientByClientID(ctx context.Context, clientID uuid.UUID) (Client, error)
}

type client struct {
	ID           uuid.UUID
	Name         string
	RedirectURIs []string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (c *client) IsNotFound() bool {
	return c.ID == uuid.Nil
}

func (c *client) IsRedirectURIMatch(redirectURI string) bool {
	for _, uri := range c.RedirectURIs {
		if uri == redirectURI {
			return true
		}
	}
	return false
}
