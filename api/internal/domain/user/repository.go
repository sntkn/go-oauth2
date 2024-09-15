package user

import (
	"log"

	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/api/internal/infrastructure/db/query"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(id string) (*User, error)
}

type Repository struct {
	query *query.Query
	gorm  *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		query: query.Use(db),
		gorm:  db,
	}
}

func (r *Repository) FindByID(id string) (*User, error) {
	userQuery := r.query.User
	user, err := userQuery.Where(userQuery.ID.Eq(id)).First()

	if err != nil {
		log.Printf("ユーザーのクエリに失敗しました: %v", err)
		return nil, err
	}

	userID, err := uuid.Parse(user.ID)
	if err != nil {
		log.Printf("UUIDの解析に失敗しました: %v", err)
		return nil, err
	}

	return &User{
		ID:        userID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		// ここに他のフィールドを追加
	}, nil
}
