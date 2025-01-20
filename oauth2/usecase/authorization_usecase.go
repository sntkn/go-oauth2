package usecase

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/domain/authorization"
	"github.com/sntkn/go-oauth2/oauth2/infrastructure/model"
	"github.com/sntkn/go-oauth2/oauth2/internal/accesstoken"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/pkg/str"
)

type IAuthorizationUsecase interface {
	Consent(uuid.UUID) (*authorization.Client, error)
	GenerateAuthorizationCode(GenerateAuthorizationCodeParams) (*authorization.AuthorizationCode, error)
	GenerateTokenByCode(string) (*authorization.Token, error)
	GenerateTokenByRefreshToken(string) (*authorization.Token, error)
	// GenerateAuthorizationCode(user *model.User, client *model.Client, scopes []string) (*model.AuthorizationCode, error)
	// ValidateAuthorizationCode(code string, clientID string) (*model.AuthorizationCode, error)
}

func NewAuthorizationUsecase(repo authorization.IAuthorizationRepository, config *config.Config, tokenGen accesstoken.Generator) IAuthorizationUsecase {
	return &AuthorizationUsecase{
		repo:     repo,
		config:   config,
		tokenGen: tokenGen,
	}
}

type AuthorizationUsecase struct {
	repo     authorization.IAuthorizationRepository
	config   *config.Config
	tokenGen accesstoken.Generator
}

func (uc *AuthorizationUsecase) Consent(clientID uuid.UUID) (*authorization.Client, error) {
	cli, err := uc.repo.FindClientByClientID(clientID)
	if err != nil {
		return nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	client := authorization.NewClient(cli.ID, cli.Name, cli.RedirectURIs, cli.CreatedAt, cli.UpdatedAt)
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

func (uc *AuthorizationUsecase) GenerateAuthorizationCode(p GenerateAuthorizationCodeParams) (*authorization.AuthorizationCode, error) {
	randomString, err := authorization.GenerateCode()
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

	code := &model.AuthorizationCode{
		Code:        randomString,
		ClientID:    clientID,
		UserID:      userID,
		Scope:       p.Scope,
		RedirectURI: p.RedirectURI,
		ExpiresAt:   time.Now().Add(time.Duration(p.Expires) * time.Second),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = uc.repo.StoreAuthorizationCode(code)
	if err != nil {
		return nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	client := authorization.NewAuthorizationCode(
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

func (uc *AuthorizationUsecase) GenerateTokenByCode(code string) (*authorization.Token, error) {
	var atokn *authorization.Token
	const (
		randomStringLen = 32
		day             = 24 * time.Hour
	)

	c, err := uc.repo.FindValidAuthorizationCode(code, time.Now())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// TODO: redirect to autorize with parameters
			return atokn, errors.NewUsecaseError(http.StatusForbidden, err.Error())
		}
		return atokn, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	currentTime := time.Now()
	if currentTime.After(c.ExpiresAt) {
		return atokn, errors.NewUsecaseError(http.StatusForbidden, "code has expired")
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
		return atokn, errors.NewUsecaseError(http.StatusInternalServerError, "code has expired")
	}

	if err = uc.repo.StoreToken(&model.Token{
		AccessToken: accessToken,
		ClientID:    c.ClientID,
		UserID:      c.UserID,
		Scope:       c.Scope,
		ExpiresAt:   expiration,
	}); err != nil {
		return atokn, errors.NewUsecaseError(http.StatusInternalServerError, "code has expired")
	}

	randomString, err := str.GenerateRandomString(randomStringLen)
	if err != nil {
		return atokn, errors.NewUsecaseError(http.StatusInternalServerError, "code has expired")
	}
	refreshExpiration := time.Now().Add(time.Duration(uc.config.AuthRefreshTokenExpiresDay) * day)
	if err = uc.repo.StoreRefreshToken(&model.RefreshToken{
		RefreshToken: randomString,
		AccessToken:  accessToken,
		ExpiresAt:    refreshExpiration,
	}); err != nil {
		return atokn, errors.NewUsecaseError(http.StatusInternalServerError, "code has expired")
	}

	// revoke code
	if err = uc.repo.RevokeCode(code); err != nil {
		return atokn, errors.NewUsecaseError(http.StatusInternalServerError, "code has expired")
	}

	atokn = &authorization.Token{
		AccessToken: accessToken,
		RefreshToken: authorization.RefreshToken{
			RefreshToken: randomString,
		},
		Expiry: expiration.Unix(),
	}

	return atokn, nil
}

func (uc *AuthorizationUsecase) GenerateTokenByRefreshToken(refreshToken string) (*authorization.Token, error) {
	return nil, nil
}
