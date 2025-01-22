package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sntkn/go-oauth2/oauth2/domain/authorization"
	"github.com/sntkn/go-oauth2/oauth2/infrastructure/model"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

func NewAuthorizationRepository(db *sqlx.DB) *AuthorizationRepository {
	return &AuthorizationRepository{
		db: db,
	}
}

type AuthorizationRepository struct {
	db *sqlx.DB
}

func (r *AuthorizationRepository) FindClientByClientID(clientID uuid.UUID) (*authorization.Client, error) {
	q := "SELECT id, redirect_uris FROM oauth2_clients WHERE id = $1"
	var c model.Client

	err := r.db.Get(&c, q, &clientID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &authorization.Client{}, nil
		}
		return nil, errors.WithStack(err)
	}

	return &authorization.Client{
		ID:           c.ID,
		RedirectURIs: c.RedirectURIs,
	}, nil
}

func (r *AuthorizationRepository) FindAuthorizationCode(code string) (*authorization.AuthorizationCode, error) {
	q := "SELECT user_id, client_id, scope, expires_at FROM oauth2_codes WHERE code = $1 AND revoked_at IS NULL"

	var c model.AuthorizationCode

	err := r.db.Get(&c, q, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &authorization.AuthorizationCode{}, nil
		}
		return nil, errors.WithStack(err)
	}
	return &authorization.AuthorizationCode{
		UserID:    c.UserID,
		ClientID:  c.ClientID,
		Scope:     c.Scope,
		ExpiresAt: c.ExpiresAt,
	}, nil
}

func (r *AuthorizationRepository) FindValidAuthorizationCode(code string, expiresAt time.Time) (*authorization.AuthorizationCode, error) {
	q := "SELECT user_id, client_id, scope, expires_at FROM oauth2_codes WHERE code = $1 AND revoked_at IS NULL AND expires_at > $2"
	var c model.AuthorizationCode

	err := r.db.Get(&c, q, code, expiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &authorization.AuthorizationCode{}, nil
		}
		return nil, errors.WithStack(err)
	}
	return &authorization.AuthorizationCode{
		UserID:    c.UserID,
		ClientID:  c.ClientID,
		Scope:     c.Scope,
		ExpiresAt: c.ExpiresAt,
	}, nil
}

func (r *AuthorizationRepository) StoreAuthorizationCode(c *authorization.AuthorizationCode) error {
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

func (r *AuthorizationRepository) StoreToken(t *authorization.Token) error {
	token := &model.Token{
		AccessToken: t.AccessToken,
		ClientID:    t.ClientID,
		UserID:      t.UserID,
		Scope:       t.Scope,
		ExpiresAt:   t.ExpiresAt,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	q := `
		INSERT INTO oauth2_tokens (access_token, client_id, user_id, scope, expires_at, created_at, updated_at)
		VALUES (:access_token, :client_id, :user_id, :scope, :expires_at, :created_at, :updated_at)
	`
	_, err := r.db.NamedExec(q, token)
	return errors.WithStack(err)
}

func (r *AuthorizationRepository) StoreRefreshToken(t *authorization.RefreshToken) error {
	rtoken := &model.RefreshToken{
		RefreshToken: t.RefreshToken,
		AccessToken:  t.AccessToken,
		ExpiresAt:    t.ExpiresAt,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	q := `INSERT INTO oauth2_refresh_tokens (refresh_token, access_token, expires_at, created_at, updated_at)
	VALUES (:refresh_token, :access_token, :expires_at, :created_at, :updated_at)`
	_, err := r.db.NamedExec(q, rtoken)
	return errors.WithStack(err)
}

func (r *AuthorizationRepository) RevokeCode(code string) error {
	updateQuery := "UPDATE oauth2_codes SET revoked_at = $1 WHERE code = $2"
	_, err := r.db.Exec(updateQuery, time.Now(), code)
	return errors.WithStack(err)
}

func (r *AuthorizationRepository) FindValidRefreshToken(refreshToken string, expiresAt time.Time) (*authorization.RefreshToken, error) {
	q := "SELECT access_token FROM oauth2_refresh_tokens WHERE refresh_token = $1 AND expires_at > $2"
	var rtkn model.RefreshToken

	err := r.db.Get(&rtkn, q, refreshToken, expiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &authorization.RefreshToken{}, nil
		}
		return nil, errors.WithStack(err)
	}

	return &authorization.RefreshToken{
		AccessToken: rtkn.AccessToken,
	}, nil
}

func (r *AuthorizationRepository) FindToken(accessToken string) (*authorization.Token, error) {
	q := "SELECT user_id, client_id, scope FROM oauth2_tokens WHERE access_token = $1"
	var tkn model.Token

	err := r.db.Get(&tkn, q, accessToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &authorization.Token{}, nil
		}
		return nil, errors.WithStack(err)
	}
	return &authorization.Token{
		UserID:   tkn.UserID,
		ClientID: tkn.ClientID,
		Scope:    tkn.Scope,
	}, nil
}

func (r *AuthorizationRepository) RevokeToken(accessToken string) error {
	updateQuery := "UPDATE oauth2_tokens SET revoked_at = $1 WHERE access_token = $2"
	_, err := r.db.Exec(updateQuery, time.Now(), accessToken)
	return errors.WithStack(err)
}

func (r *AuthorizationRepository) RevokeRefreshToken(refreshToken string) error {
	updateQuery := "UPDATE oauth2_refresh_tokens SET revoked_at = $1 WHERE refresh_token = $2"
	_, err := r.db.Exec(updateQuery, time.Now(), refreshToken)
	return errors.WithStack(err)
}
