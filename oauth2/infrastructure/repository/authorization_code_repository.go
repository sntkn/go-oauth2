package repository

import (
	"context"
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

func (r *AuthorizationCodeRepository) FindAuthorizationCode(ctx context.Context, code string) (domain.AuthorizationCode, error) {
	q := "SELECT user_id, client_id, scope, expires_at FROM oauth2_codes WHERE code = $1 AND revoked_at IS NULL"

	var c model.AuthorizationCode

	err := r.db.GetContext(ctx, &c, q, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}
	return domain.NewAuthorizationCode(domain.AuthorizationCodeParams{
		UserID:    c.UserID.String(),
		ClientID:  c.ClientID.String(),
		Scope:     c.Scope,
		ExpiresAt: c.ExpiresAt,
	})
}

func (r *AuthorizationCodeRepository) FindValidAuthorizationCode(
	ctx context.Context,
	code string,
	expiresAt time.Time,
) (domain.AuthorizationCode, error) {
	q := "SELECT user_id, client_id, scope, expires_at FROM oauth2_codes WHERE code = $1 AND revoked_at IS NULL AND expires_at > $2"
	var c model.AuthorizationCode

	err := r.db.GetContext(ctx, &c, q, code, expiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}
	return domain.NewAuthorizationCode(domain.AuthorizationCodeParams{
		UserID:    c.UserID.String(),
		ClientID:  c.ClientID.String(),
		Scope:     c.Scope,
		ExpiresAt: c.ExpiresAt,
	})
}

func (r *AuthorizationCodeRepository) StoreAuthorizationCode(ctx context.Context, p domain.StoreAuthorizationCodeParams) (string, error) {
	m := &model.AuthorizationCode{
		Code:        p.Code,
		ClientID:    p.ClientID,
		UserID:      p.UserID,
		Scope:       p.Scope,
		RedirectURI: p.RedirectURI,
		ExpiresAt:   p.ExpiresAt,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	q := `
			INSERT INTO oauth2_codes
				(code, client_id, user_id, scope, redirect_uri, expires_at, created_at, updated_at)
			VALUES
				(:code, :client_id, :user_id, :scope, :redirect_uri, :expires_at, :created_at, :updated_at)
	`

	_, err := r.db.NamedExecContext(ctx, q, m)
	if err != nil {
		return "", errors.WithStack(err)
	}

	// returning primary key
	return p.Code, nil
}

func (r *AuthorizationCodeRepository) RevokeCode(ctx context.Context, code string) error {
	updateQuery := "UPDATE oauth2_codes SET revoked_at = $1 WHERE code = $2"
	_, err := r.db.ExecContext(ctx, updateQuery, time.Now(), code)
	return errors.WithStack(err)
}
