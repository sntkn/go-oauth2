package repository

import (
	"log"

	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/api/internal/domain/user"
	"github.com/sntkn/go-oauth2/api/internal/infrastructure/db/query"
	"gorm.io/gorm"
)

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

func (r *Repository) FindByID(id string) (*user.User, error) {
	userQuery := r.query.User
	u, err := userQuery.Where(userQuery.ID.Eq(id)).First()

	if err != nil {
		log.Printf("ユーザーのクエリに失敗しました: %v", err)
		return nil, err
	}

	userID, err := uuid.Parse(u.ID)
	if err != nil {
		log.Printf("UUIDの解析に失敗しました: %v", err)
		return nil, err
	}

	return &user.User{
		ID:        userID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		// ここに他のフィールドを追加
	}, nil
}
