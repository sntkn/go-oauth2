package domain

import (
	"time"

	"github.com/google/uuid"
)

type ClientRepository interface {
	FindClientByClientID(clientID uuid.UUID) (*Client, error)
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
