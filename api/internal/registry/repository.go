package registry

import (
	"log"

	"github.com/sntkn/go-oauth2/api/internal/infrastructure/db/model"
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

func (r *Repository) FindByID(id string) (*model.User, error) {
	userQuery := r.query.User
	u, err := userQuery.Where(userQuery.ID.Eq(id)).First()

	if err != nil {
		log.Printf("ユーザーのクエリに失敗しました: %v", err)
		return nil, err
	}

	return u, nil
}
