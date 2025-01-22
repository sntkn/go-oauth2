package repository

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sntkn/go-oauth2/oauth2/domain/authentication"
	"github.com/sntkn/go-oauth2/oauth2/infrastructure/model"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

func NewAuthenticationRepository(db *sqlx.DB) *AuthenticationRepository {
	return &AuthenticationRepository{
		db: db,
	}
}

type AuthenticationRepository struct {
	db *sqlx.DB
}

func (r *AuthenticationRepository) FindClientByClientID(clientID uuid.UUID) (*authentication.Client, error) {
	q := "SELECT id, redirect_uris FROM oauth2_clients WHERE id = $1"
	var c model.Client

	err := r.db.Get(&c, q, &clientID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &authentication.Client{}, nil
		}
		return nil, errors.WithStack(err)
	}

	return &authentication.Client{
		ID:           c.ID,
		RedirectURIs: c.RedirectURIs,
	}, nil
}

func (r *AuthenticationRepository) FindUserByEmail(email string) (*authentication.User, error) {
	q := "SELECT id, email, password FROM users WHERE email = $1"
	var user model.User

	err := r.db.Get(&user, q, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &authentication.User{}, nil
		}
		return nil, errors.WithStack(err)
	}
	return &authentication.User{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}
