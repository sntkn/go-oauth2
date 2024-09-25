package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/api/internal/infrastructure/db/model"
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

func NewTimeline(posts []*model.Post) ([]*Timeline, error) {
	var tl []*Timeline

	for _, l := range posts {
		tl = append(tl, &Timeline{
			PostTime:    l.CreatedAt,
			Content:     l.Content,
			PostUser:    User{},
			Inpressions: 0,
			Likes:       []UserID{},
			Reposts:     []Timeline{},
		})
	}

	return tl, nil
}
