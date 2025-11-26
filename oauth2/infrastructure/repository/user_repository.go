package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/sntkn/go-oauth2/oauth2/domain"
	"github.com/sntkn/go-oauth2/oauth2/infrastructure/model"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
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
	var u model.User

	err := r.db.GetContext(ctx, &u, q, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.NewUser(domain.UserParams{}), nil
		}
		return nil, errors.WithStack(err)
	}

	return domain.NewUser(domain.UserParams{
		ID:       u.ID,
		Email:    u.Email,
		Password: u.Password,
	}), nil
}
