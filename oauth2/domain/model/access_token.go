package model

import (
	"time"

	"github.com/google/uuid"
)

type AccessToken struct {
	AccessToken string
	ClientID    uuid.UUID
	UserID      uuid.UUID
	Scope       string
	ExpiresAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
