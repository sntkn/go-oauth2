package model

import (
	"time"

	"github.com/google/uuid"
)

type Client struct {
	ID           uuid.UUID
	Name         string
	RedirectURIs string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
