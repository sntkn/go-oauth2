package model

import (
	"time"

	"github.com/google/uuid"
)

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
