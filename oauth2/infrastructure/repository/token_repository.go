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

func NewTokenRepository(db *sqlx.DB) *TokenRepository {
	return &TokenRepository{
		db: db,
	}
}

type TokenRepository struct {
	db *sqlx.DB
}

func (r *TokenRepository) StoreToken(ctx context.Context, accessToken domain.Token) error {
	m := &model.Token{
		AccessToken: accessToken.GetAccessToken(),
		ClientID:    accessToken.GetClientID(),
		UserID:      accessToken.GetUserID(),
		Scope:       accessToken.GetScope(),
		ExpiresAt:   accessToken.GetExpiresAt(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	q := `
		INSERT INTO oauth2_tokens (access_token, client_id, user_id, scope, expires_at, created_at, updated_at)
		VALUES (:access_token, :client_id, :user_id, :scope, :expires_at, :created_at, :updated_at)
	`
	_, err := r.db.NamedExecContext(ctx, q, m)
	return errors.WithStack(err)
}

func (r *TokenRepository) FindToken(ctx context.Context, accessToken string) (domain.Token, error) {
	q := "SELECT user_id, client_id, scope FROM oauth2_tokens WHERE access_token = $1"
	var tkn model.Token

	err := r.db.GetContext(ctx, &tkn, q, accessToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}
	return domain.NewToken(domain.TokenParams{
		UserID:   tkn.UserID,
		ClientID: tkn.ClientID,
		Scope:    tkn.Scope,
	}), nil
}

func (r *TokenRepository) RevokeToken(ctx context.Context, accessToken string) error {
	updateQuery := "UPDATE oauth2_tokens SET revoked_at = $1 WHERE access_token = $2"
	_, err := r.db.ExecContext(ctx, updateQuery, time.Now(), accessToken)
	return errors.WithStack(err)
}
