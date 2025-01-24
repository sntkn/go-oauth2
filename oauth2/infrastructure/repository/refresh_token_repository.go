package repository

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sntkn/go-oauth2/oauth2/domain/refresh_token"
	"github.com/sntkn/go-oauth2/oauth2/infrastructure/model"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

func NewRefreshTokenRepository(db *sqlx.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		db: db,
	}
}

type RefreshTokenRepository struct {
	db *sqlx.DB
}

func (r *RefreshTokenRepository) StoreRefreshToken(t *refresh_token.RefreshToken) error {
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

func (r *RefreshTokenRepository) FindValidRefreshToken(refreshToken string, expiresAt time.Time) (*refresh_token.RefreshToken, error) {
	q := "SELECT access_token FROM oauth2_refresh_tokens WHERE refresh_token = $1 AND expires_at > $2"
	var rtkn model.RefreshToken

	err := r.db.Get(&rtkn, q, refreshToken, expiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &refresh_token.RefreshToken{}, nil
		}
		return nil, errors.WithStack(err)
	}

	return &refresh_token.RefreshToken{
		AccessToken: rtkn.AccessToken,
	}, nil
}

func (r *RefreshTokenRepository) RevokeRefreshToken(refreshToken string) error {
	updateQuery := "UPDATE oauth2_refresh_tokens SET revoked_at = $1 WHERE refresh_token = $2"
	_, err := r.db.Exec(updateQuery, time.Now(), refreshToken)
	return errors.WithStack(err)
}
