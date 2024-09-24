package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserID uuid.UUID

type Timeline struct {
	PostTime    time.Time
	Content     string
	PostUser    User
	Inpressions uint32
	Likes       []UserID
	Reposts     []Timeline
}

type User struct {
	ID   UserID
	Name string
	// IconURI string
}
