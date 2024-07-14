package usecases

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/accesstoken"
	"github.com/sntkn/go-oauth2/oauth2/internal/entity"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/pkg/str"
)

type CreateTokenByCode struct {
	cfg *config.Config
	db  *repository.Repository
}

func NewCreateTokenByCode(cfg *config.Config, db *repository.Repository) *CreateTokenByCode {
	return &CreateTokenByCode{
		cfg: cfg,
		db:  db,
	}
}

func (u *CreateTokenByCode) Invoke(c *gin.Context, authCode string) (entity.AuthTokens, error) {
	var atokn entity.AuthTokens
	const (
		randomStringLen = 32
		day             = 24 * time.Hour
	)

	code, err := u.db.FindValidOAuth2Code(authCode, time.Now())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// TODO: redirect to autorize with parameters
			return atokn, cerrs.NewUsecaseError(http.StatusBadRequest, err.Error())
		}
		return atokn, cerrs.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	currentTime := time.Now()
	if currentTime.After(code.ExpiresAt) {
		return atokn, cerrs.NewUsecaseError(http.StatusForbidden, "code has expired")
	}

	// create token and refresh token
	expiration := time.Now().Add(u.cfg.AuthTokenExpiresMin * time.Minute)
	t := accesstoken.TokenParams{
		UserID:    code.UserID,
		ClientID:  code.ClientID,
		Scope:     code.Scope,
		ExpiresAt: expiration,
	}
	accessToken, err := accesstoken.Generate(t)
	if err != nil {
		return atokn, cerrs.NewUsecaseError(http.StatusInternalServerError, "code has expired")
	}

	if err = u.db.RegisterToken(&repository.Token{
		AccessToken: accessToken,
		ClientID:    code.ClientID,
		UserID:      code.UserID,
		Scope:       code.Scope,
		ExpiresAt:   expiration,
	}); err != nil {
		return atokn, cerrs.NewUsecaseError(http.StatusInternalServerError, "code has expired")
	}

	randomString, err := str.GenerateRandomString(randomStringLen)
	if err != nil {
		return atokn, cerrs.NewUsecaseError(http.StatusInternalServerError, "code has expired")
	}
	refreshExpiration := time.Now().Add(u.cfg.AuthRefreshTokenExpiresDay * day)
	if err = u.db.RegesterRefreshToken(&repository.RefreshToken{
		RefreshToken: randomString,
		AccessToken:  accessToken,
		ExpiresAt:    refreshExpiration,
	}); err != nil {
		return atokn, cerrs.NewUsecaseError(http.StatusInternalServerError, "code has expired")
	}

	// revoke code
	if err = u.db.RevokeCode(authCode); err != nil {
		return atokn, cerrs.NewUsecaseError(http.StatusInternalServerError, "code has expired")
	}

	return entity.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: randomString,
		Expiry:       expiration.Unix(),
	}, nil
}
