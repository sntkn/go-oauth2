package usecases

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/sntkn/go-oauth2/oauth2/internal/accesstoken"
	"github.com/sntkn/go-oauth2/oauth2/internal/entity"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/pkg/str"
)

type CreateTokenByCode struct {
	cfg      *config.Config
	db       repository.OAuth2Repository
	tokenGen accesstoken.Generator
}

func NewCreateTokenByCode(cfg *config.Config, db repository.OAuth2Repository, tokenGen accesstoken.Generator) *CreateTokenByCode {
	return &CreateTokenByCode{
		cfg:      cfg,
		db:       db,
		tokenGen: tokenGen,
	}
}

func (u *CreateTokenByCode) Invoke(authCode string) (*entity.AuthTokens, error) {
	var atokn *entity.AuthTokens
	const (
		randomStringLen = 32
		day             = 24 * time.Hour
	)

	code, err := u.db.FindValidOAuth2Code(authCode, time.Now())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// TODO: redirect to autorize with parameters
			return atokn, errors.NewUsecaseError(http.StatusForbidden, err.Error())
		}
		return atokn, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	currentTime := time.Now()
	if currentTime.After(code.ExpiresAt) {
		return atokn, errors.NewUsecaseError(http.StatusForbidden, "code has expired")
	}

	// create token and refresh token
	expiration := time.Now().Add(time.Duration(u.cfg.AuthTokenExpiresMin) * time.Minute)
	t := accesstoken.TokenParams{
		UserID:    code.UserID,
		ClientID:  code.ClientID,
		Scope:     code.Scope,
		ExpiresAt: expiration,
	}
	accessToken, err := u.tokenGen.Generate(&t, u.cfg.PrivateKey)
	if err != nil {
		return atokn, errors.NewUsecaseError(http.StatusInternalServerError, "code has expired")
	}

	if err = u.db.RegisterToken(&repository.Token{
		AccessToken: accessToken,
		ClientID:    code.ClientID,
		UserID:      code.UserID,
		Scope:       code.Scope,
		ExpiresAt:   expiration,
	}); err != nil {
		return atokn, errors.NewUsecaseError(http.StatusInternalServerError, "code has expired")
	}

	randomString, err := str.GenerateRandomString(randomStringLen)
	if err != nil {
		return atokn, errors.NewUsecaseError(http.StatusInternalServerError, "code has expired")
	}
	refreshExpiration := time.Now().Add(time.Duration(u.cfg.AuthRefreshTokenExpiresDay) * day)
	if err = u.db.RegisterRefreshToken(&repository.RefreshToken{
		RefreshToken: randomString,
		AccessToken:  accessToken,
		ExpiresAt:    refreshExpiration,
	}); err != nil {
		return atokn, errors.NewUsecaseError(http.StatusInternalServerError, "code has expired")
	}

	// revoke code
	if err = u.db.RevokeCode(authCode); err != nil {
		return atokn, errors.NewUsecaseError(http.StatusInternalServerError, "code has expired")
	}

	return &entity.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: randomString,
		Expiry:       expiration.Unix(),
	}, nil
}
