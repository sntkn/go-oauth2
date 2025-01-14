package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sntkn/go-oauth2/oauth2/infrastructure/model"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

func NewDB(host, user, password, dbName string, port int) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	// PostgreSQLに接続
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return db, nil
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

type Repository struct {
	db *sqlx.DB
}

func (r *Repository) FindClientByClientID(clientID uuid.UUID) (*model.Client, error) {
	q := "SELECT id, redirect_uris FROM oauth2_clients WHERE id = $1"
	var c model.Client

	err := r.db.Get(&c, q, &clientID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &model.Client{}, nil
		}
		return nil, errors.WithStack(err)
	}
	return &c, errors.WithStack(err)
}

func (r *Repository) FindUserByEmail(email string) (*model.User, error) {
	q := "SELECT id, email, password FROM users WHERE email = $1"
	var user model.User

	err := r.db.Get(&user, q, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &model.User{}, nil
		}
		return nil, errors.WithStack(err)
	}
	return &user, errors.WithStack(err)
}

func (r *Repository) FindAuthorizationCode(code string) (*model.AuthorizationCode, error) {
	q := "SELECT user_id, client_id, scope, expires_at FROM oauth2_codes WHERE code = $1 AND revoked_at IS NULL"

	var c model.AuthorizationCode

	err := r.db.Get(&c, q, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &model.AuthorizationCode{}, nil
		}
		return nil, errors.WithStack(err)
	}
	return &c, errors.WithStack(err)
}

func (r *Repository) StoreAuthorizationCode(c *model.AuthorizationCode) error {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	q := `
		INSERT INTO oauth2_refresh_tokens
			(refresh_token, access_token, expires_at, created_at, updated_at)
		VALUES
			(:refresh_token, :access_token, :expires_at, :created_at, :updated_at)
	`

	_, err := r.db.NamedExec(q, c)
	return errors.WithStack(err)
}
