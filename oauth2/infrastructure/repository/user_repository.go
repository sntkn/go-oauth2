package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/sntkn/go-oauth2/oauth2/domain/model"
)

func NewUserRepository(db *sqlx.DB) IUserRepository {
	return UserRepository{
		DB: db,
	}
}

type IUserRepository interface {
	FindUserByEmail(email string) (*model.User, error)
}

type UserRepository struct {
	DB *sqlx.DB
}

func (r UserRepository) FindUserByEmail(email string) (*model.User, error) {
	return nil, nil
}
