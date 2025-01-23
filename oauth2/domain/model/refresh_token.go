package model

import "time"

type RefreshToken struct {
	RefreshToken string
	AccessToken  string
	ExpiresAt    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
