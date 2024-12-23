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

type CreateTokenByRefreshToken struct {
	cfg      *config.Config
	db       repository.OAuth2Repository
	tokenGen accesstoken.Generator
}

func NewCreateTokenByRefreshToken(cfg *config.Config, db repository.OAuth2Repository, tokenGen accesstoken.Generator) *CreateTokenByRefreshToken {
	return &CreateTokenByRefreshToken{
		cfg:      cfg,
		db:       db,
		tokenGen: tokenGen,
	}
}

func (u *CreateTokenByRefreshToken) Invoke(refreshToken string) (*entity.AuthTokens, error) {
	var atokn *entity.AuthTokens
	const (
		randomStringLen = 32
		day             = 24 * time.Hour
	)

	// TODO: find refresh token, if not expired
	rt, err := u.db.FindValidRefreshToken(refreshToken, time.Now())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return atokn, errors.NewUsecaseError(http.StatusForbidden, err.Error())
		}
		return atokn, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}
	// find access token

	// TODO: create token and refresh token
	tkn, err := u.db.FindToken(rt.AccessToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// TODO: redirect to autorize with parameters
			return atokn, errors.NewUsecaseError(http.StatusForbidden, err.Error())
		}
		return atokn, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}
	expiration := time.Now().Add(time.Duration(u.cfg.AuthTokenExpiresMin) * time.Minute)

	t := &accesstoken.TokenParams{
		UserID:    tkn.UserID,
		ClientID:  tkn.ClientID,
		Scope:     tkn.Scope,
		ExpiresAt: expiration,
	}
	accessToken, err := u.tokenGen.Generate(t, u.cfg.PrivateKey)

	if err != nil {
		return atokn, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	if err = u.db.RegisterToken(&repository.Token{
		AccessToken: accessToken,
		ClientID:    tkn.ClientID,
		UserID:      tkn.UserID,
		Scope:       tkn.Scope,
		ExpiresAt:   expiration,
	}); err != nil {
		return atokn, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	randomString, err := str.GenerateRandomString(randomStringLen)
	if err != nil {
		return atokn, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}
	refreshExpiration := time.Now().Add(time.Duration(u.cfg.AuthRefreshTokenExpiresDay) * day)

	if err = u.db.RegisterRefreshToken(&repository.RefreshToken{
		RefreshToken: randomString,
		AccessToken:  accessToken,
		ExpiresAt:    refreshExpiration,
	}); err != nil {
		return atokn, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	// TODO: revoke old token and refresh token
	if err = u.db.RevokeToken(tkn.AccessToken); err != nil {
		return atokn, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}
	if err = u.db.RevokeRefreshToken(rt.RefreshToken); err != nil {
		return atokn, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	return &entity.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: randomString,
		Expiry:       expiration.Unix(),
	}, nil
}
