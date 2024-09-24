package domain

import (
	"log"
	"time"

	"github.com/google/uuid"
	model "github.com/sntkn/go-oauth2/api/internal/infrastructure/db/model"
)

type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(u *model.User) (*User, error) {
	userID, err := uuid.Parse(u.ID)
	if err != nil {
		log.Printf("UUIDの解析に失敗しました: %v", err)
		return nil, err
	}

	return &User{
		ID:        userID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		// ここに他のフィールドを追加
	}, nil
}
