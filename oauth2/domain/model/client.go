package model

import (
	"time"

	"github.com/google/uuid"
)

func NewClient(id uuid.UUID, name, redirectURIs string, createdAt, updatedAt time.Time) *Client {
	return &Client{
		ID:           id,
		Name:         name,
		RedirectURIs: redirectURIs,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}

type Client struct {
	ID           uuid.UUID
	Name         string
	RedirectURIs string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
