package usecase

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/domain/authorization_code"
	"github.com/sntkn/go-oauth2/oauth2/domain/client"
	"github.com/sntkn/go-oauth2/oauth2/domain/refresh_token"
	"github.com/sntkn/go-oauth2/oauth2/domain/token"
	"github.com/sntkn/go-oauth2/oauth2/internal/common/accesstoken"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/pkg/str"
)

type IAuthorizationUsecase interface {
	Consent(uuid.UUID) (*client.Client, error)
	GenerateAuthorizationCode(GenerateAuthorizationCodeParams) (*authorization_code.AuthorizationCode, error)
	GenerateTokenByCode(string) (*token.Token, *refresh_token.RefreshToken, error)
	GenerateTokenByRefreshToken(string) (*token.Token, *refresh_token.RefreshToken, error)
	// GenerateAuthorizationCode(user *model.User, client *model.Client, scopes []string) (*model.AuthorizationCode, error)
	// ValidateAuthorizationCode(code string, clientID string) (*model.AuthorizationCode, error)
}

func NewAuthorizationUsecase(
	clientRepo client.ClientRepository,
	codeRepo authorization_code.AuthorizationCodeRepository,
	tokenRepo token.TokenRepository,
	refreshTokenRepo refresh_token.RefreshTokenRepository,

	config *config.Config, tokenGen accesstoken.Generator) IAuthorizationUsecase {
	return &AuthorizationUsecase{
		clientRepo:       clientRepo,
		codeRepo:         codeRepo,
		tokenRepo:        tokenRepo,
		refreshTokenRepo: refreshTokenRepo,
		config:           config,
		tokenGen:         tokenGen,
	}
}

type AuthorizationUsecase struct {
	clientRepo       client.ClientRepository
	codeRepo         authorization_code.AuthorizationCodeRepository
	tokenRepo        token.TokenRepository
	refreshTokenRepo refresh_token.RefreshTokenRepository
	config           *config.Config
	tokenGen         accesstoken.Generator
}

func (uc *AuthorizationUsecase) Consent(clientID uuid.UUID) (*client.Client, error) {
	cli, err := uc.clientRepo.FindClientByClientID(clientID)
	if err != nil {
		return nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	client := client.NewClient(cli.ID, cli.Name, cli.RedirectURIs, cli.CreatedAt, cli.UpdatedAt)
	if client.IsNotFound() {
		return nil, errors.NewUsecaseError(http.StatusBadRequest, "client not found")
	}

	return client, nil
}

type GenerateAuthorizationCodeParams struct {
	UserID      string
	ClientID    string
	RedirectURI string
	Scope       string
	Expires     int
}

func (uc *AuthorizationUsecase) GenerateAuthorizationCode(p GenerateAuthorizationCodeParams) (*authorization_code.AuthorizationCode, error) {
	randomString, err := authorization_code.GenerateCode()
	if err != nil {
		return nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	clientID, err := uuid.Parse(p.ClientID)
	if err != nil {
		return nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}
	userID, err := uuid.Parse(p.UserID)
	if err != nil {
		return nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	code := &authorization_code.AuthorizationCode{
		Code:        randomString,
		ClientID:    clientID,
		UserID:      userID,
		Scope:       p.Scope,
		RedirectURI: p.RedirectURI,
		ExpiresAt:   time.Now().Add(time.Duration(p.Expires) * time.Second),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = uc.codeRepo.StoreAuthorizationCode(code)
	if err != nil {
		return nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	client := authorization_code.NewAuthorizationCode(
		code.Code,
		code.ClientID,
		code.UserID,
		code.Scope,
		code.RedirectURI,
		code.ExpiresAt,
		code.CreatedAt,
		code.UpdatedAt,
	)

	return client, nil
}

func (uc *AuthorizationUsecase) GenerateTokenByCode(code string) (*token.Token, *refresh_token.RefreshToken, error) {
	var atokn *token.Token
	var rtokn *refresh_token.RefreshToken
	const (
		randomStringLen = 32
		day             = 24 * time.Hour
	)

	c, err := uc.codeRepo.FindValidAuthorizationCode(code, time.Now())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// TODO: redirect to autorize with parameters
			return atokn, rtokn, errors.NewUsecaseError(http.StatusForbidden, err.Error())
		}
		return atokn, rtokn, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	currentTime := time.Now()
	if currentTime.After(c.ExpiresAt) {
		return atokn, rtokn, errors.NewUsecaseError(http.StatusForbidden, "code has expired")
	}

	// create token and refresh token
	expiration := time.Now().Add(time.Duration(uc.config.AuthTokenExpiresMin) * time.Minute)
	t := accesstoken.TokenParams{
		UserID:    c.UserID,
		ClientID:  c.ClientID,
		Scope:     c.Scope,
		ExpiresAt: expiration,
	}
	accessToken, err := uc.tokenGen.Generate(&t, uc.config.PrivateKey)
	if err != nil {
		return atokn, rtokn, errors.NewUsecaseError(http.StatusInternalServerError, "code has expired")
	}

	if err = uc.tokenRepo.StoreToken(&token.Token{
		AccessToken: accessToken,
		ClientID:    c.ClientID,
		UserID:      c.UserID,
		Scope:       c.Scope,
		ExpiresAt:   expiration,
	}); err != nil {
		return atokn, rtokn, errors.NewUsecaseError(http.StatusInternalServerError, "code has expired")
	}

	randomString, err := str.GenerateRandomString(randomStringLen)
	if err != nil {
		return atokn, rtokn, errors.NewUsecaseError(http.StatusInternalServerError, "code has expired")
	}
	refreshExpiration := time.Now().Add(time.Duration(uc.config.AuthRefreshTokenExpiresDay) * day)
	if err = uc.refreshTokenRepo.StoreRefreshToken(&refresh_token.RefreshToken{
		RefreshToken: randomString,
		AccessToken:  accessToken,
		ExpiresAt:    refreshExpiration,
	}); err != nil {
		return atokn, rtokn, errors.NewUsecaseError(http.StatusInternalServerError, "code has expired")
	}

	// revoke code
	if err = uc.codeRepo.RevokeCode(code); err != nil {
		return atokn, rtokn, errors.NewUsecaseError(http.StatusInternalServerError, "code has expired")
	}

	atokn = &token.Token{
		AccessToken: accessToken,
		ExpiresAt:   expiration,
	}

	rtokn = &refresh_token.RefreshToken{
		RefreshToken: randomString,
		ExpiresAt:    refreshExpiration,
	}

	return atokn, rtokn, nil
}

func (uc *AuthorizationUsecase) GenerateTokenByRefreshToken(refreshToken string) (*token.Token, *refresh_token.RefreshToken, error) {
	var atokn *token.Token
	var rtokn *refresh_token.RefreshToken
	const (
		randomStringLen = 32
		day             = 24 * time.Hour
	)

	rt, err := uc.refreshTokenRepo.FindValidRefreshToken(refreshToken, time.Now())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return atokn, rtokn, errors.NewUsecaseError(http.StatusForbidden, err.Error())
		}
		return atokn, rtokn, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}
	// find access token

	tkn, err := uc.tokenRepo.FindToken(rt.AccessToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return atokn, rtokn, errors.NewUsecaseError(http.StatusForbidden, err.Error())
		}
		return atokn, rtokn, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}
	expiration := time.Now().Add(time.Duration(uc.config.AuthTokenExpiresMin) * time.Minute)

	t := &accesstoken.TokenParams{
		UserID:    tkn.UserID,
		ClientID:  tkn.ClientID,
		Scope:     tkn.Scope,
		ExpiresAt: expiration,
	}
	accessToken, err := uc.tokenGen.Generate(t, uc.config.PrivateKey)

	if err != nil {
		return atokn, rtokn, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	if err = uc.tokenRepo.StoreToken(&token.Token{
		AccessToken: accessToken,
		ClientID:    tkn.ClientID,
		UserID:      tkn.UserID,
		Scope:       tkn.Scope,
		ExpiresAt:   expiration,
	}); err != nil {
		return atokn, rtokn, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	randomString, err := str.GenerateRandomString(randomStringLen)
	if err != nil {
		return atokn, rtokn, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}
	refreshExpiration := time.Now().Add(time.Duration(uc.config.AuthRefreshTokenExpiresDay) * day)

	if err = uc.refreshTokenRepo.StoreRefreshToken(&refresh_token.RefreshToken{
		RefreshToken: randomString,
		AccessToken:  accessToken,
		ExpiresAt:    refreshExpiration,
	}); err != nil {
		return atokn, rtokn, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	if err = uc.tokenRepo.RevokeToken(tkn.AccessToken); err != nil {
		return atokn, rtokn, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}
	if err = uc.refreshTokenRepo.RevokeRefreshToken(rt.RefreshToken); err != nil {
		return atokn, rtokn, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	atokn = &token.Token{
		AccessToken: accessToken,
		ExpiresAt:   expiration,
	}
	rtokn = &refresh_token.RefreshToken{
		RefreshToken: randomString,
		ExpiresAt:    refreshExpiration,
	}

	return atokn, rtokn, nil
}
