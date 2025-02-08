package repository

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sntkn/go-oauth2/oauth2/domain"
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

func (r *RefreshTokenRepository) StoreRefreshToken(t domain.RefreshToken) error {
	rtoken := &model.RefreshToken{
		RefreshToken: t.GetRefreshToken(),
		AccessToken:  t.GetAccessToken(),
		ExpiresAt:    t.GetExpiresAt(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	q := `INSERT INTO oauth2_refresh_tokens (domain, access_token, expires_at, created_at, updated_at)
	VALUES (:domain, :access_token, :expires_at, :created_at, :updated_at)`
	_, err := r.db.NamedExec(q, rtoken)
	return errors.WithStack(err)
}

func (r *RefreshTokenRepository) FindRefreshToken(refreshToken string) (domain.RefreshToken, error) {
	q := "SELECT access_token FROM oauth2_refresh_tokens WHERE domain = $1"
	var rtkn model.RefreshToken

	err := r.db.Get(&rtkn, q, refreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.NewRefreshToken(domain.RefreshTokenParams{}), nil
		}
		return nil, errors.WithStack(err)
	}

	return domain.NewRefreshToken(domain.RefreshTokenParams{
		AccessToken: rtkn.AccessToken,
	}), nil
}

func (r *RefreshTokenRepository) RevokeRefreshToken(refreshToken string) error {
	updateQuery := "UPDATE oauth2_refresh_tokens SET revoked_at = $1 WHERE domain = $2"
	_, err := r.db.Exec(updateQuery, time.Now(), refreshToken)
	return errors.WithStack(err)
}
