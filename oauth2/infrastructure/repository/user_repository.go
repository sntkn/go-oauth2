package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sntkn/go-oauth2/oauth2/domain"
	"github.com/sntkn/go-oauth2/oauth2/infrastructure/model"
)

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

type UserRepository struct {
	db *sqlx.DB
}

func (r *UserRepository) FindUserByEmail(ctx context.Context, email string) (domain.User, error) {
	q := "SELECT id, email, password FROM users WHERE email = $1"
	mapper := func(u model.User) (domain.User, error) {
		return domain.NewUser(domain.UserParams{
			ID:       u.ID,
			Email:    u.Email,
			Password: u.Password,
		}), nil
	}

	user, ok, err := fetchAndMap[model.User, domain.User](ctx, r.db, q, mapper, email)
	if err != nil {
		return nil, err
	}
	if !ok {
		return domain.NewUser(domain.UserParams{}), nil
	}

	return user, nil
}
