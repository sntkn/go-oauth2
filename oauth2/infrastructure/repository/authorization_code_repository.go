package repository

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sntkn/go-oauth2/oauth2/domain"
	"github.com/sntkn/go-oauth2/oauth2/infrastructure/model"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

func NewAuthorizationCodeRepository(db *sqlx.DB) *AuthorizationCodeRepository {
	return &AuthorizationCodeRepository{
		db: db,
	}
}

type AuthorizationCodeRepository struct {
	db *sqlx.DB
}

func (r *AuthorizationCodeRepository) FindAuthorizationCode(code string) (*domain.AuthorizationCode, error) {
	q := "SELECT user_id, client_id, scope, expires_at FROM oauth2_codes WHERE code = $1 AND revoked_at IS NULL"

	var c model.AuthorizationCode

	err := r.db.Get(&c, q, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &domain.AuthorizationCode{}, nil
		}
		return nil, errors.WithStack(err)
	}
	return &domain.AuthorizationCode{
		UserID:    c.UserID,
		ClientID:  c.ClientID,
		Scope:     c.Scope,
		ExpiresAt: c.ExpiresAt,
	}, nil
}

func (r *AuthorizationCodeRepository) FindValidAuthorizationCode(code string, expiresAt time.Time) (*domain.AuthorizationCode, error) {
	q := "SELECT user_id, client_id, scope, expires_at FROM oauth2_codes WHERE code = $1 AND revoked_at IS NULL AND expires_at > $2"
	var c model.AuthorizationCode

	err := r.db.Get(&c, q, code, expiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &domain.AuthorizationCode{}, nil
		}
		return nil, errors.WithStack(err)
	}
	return &domain.AuthorizationCode{
		UserID:    c.UserID,
		ClientID:  c.ClientID,
		Scope:     c.Scope,
		ExpiresAt: c.ExpiresAt,
	}, nil
}

func (r *AuthorizationCodeRepository) StoreAuthorizationCode(c *domain.AuthorizationCode) error {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	q := `
			INSERT INTO oauth2_codes
				(code, client_id, user_id, scope, redirect_uri, expires_at, created_at, updated_at)
			VALUES
				(:code, :client_id, :user_id, :scope, :redirect_uri, :expires_at, :created_at, :updated_at)
	`

	_, err := r.db.NamedExec(q, c)
	return errors.WithStack(err)
}

func (r *AuthorizationCodeRepository) RevokeCode(code string) error {
	updateQuery := "UPDATE oauth2_codes SET revoked_at = $1 WHERE code = $2"
	_, err := r.db.Exec(updateQuery, time.Now(), code)
	return errors.WithStack(err)
}
