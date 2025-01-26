package usecase

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/domain"
	"github.com/sntkn/go-oauth2/oauth2/internal/common/accesstoken"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/pkg/str"
)

type IAuthorizationUsecase interface {
	Consent(uuid.UUID) (domain.Client, error)
	GenerateAuthorizationCode(GenerateAuthorizationCodeParams) (domain.AuthorizationCode, error)
	GenerateTokenByCode(string) (*domain.Token, *domain.RefreshToken, error)
	GenerateTokenByRefreshToken(string) (*domain.Token, *domain.RefreshToken, error)
	// GenerateAuthorizationCode(user *model.User, client *model.Client, scopes []string) (*model.AuthorizationCode, error)
	// ValidateAuthorizationCode(code string, clientID string) (*model.AuthorizationCode, error)
}

func NewAuthorizationUsecase(
	clientRepo domain.ClientRepository,
	codeRepo domain.AuthorizationCodeRepository,
	tokenRepo domain.TokenRepository,
	refreshTokenRepo domain.RefreshTokenRepository,

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
	clientRepo       domain.ClientRepository
	codeRepo         domain.AuthorizationCodeRepository
	tokenRepo        domain.TokenRepository
	refreshTokenRepo domain.RefreshTokenRepository
	config           *config.Config
	tokenGen         accesstoken.Generator
}

func (uc *AuthorizationUsecase) Consent(clientID uuid.UUID) (domain.Client, error) {
	client, err := uc.clientRepo.FindClientByClientID(clientID)
	if err != nil {
		return nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

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

func (uc *AuthorizationUsecase) GenerateAuthorizationCode(p GenerateAuthorizationCodeParams) (domain.AuthorizationCode, error) {
	randomString, err := domain.GenerateCode()
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

	code := domain.NewAuthorizationCode(domain.AuthorizationCodeParams{
		Code:        randomString,
		ClientID:    clientID,
		UserID:      userID,
		Scope:       p.Scope,
		RedirectURI: p.RedirectURI,
		ExpiresAt:   time.Now().Add(time.Duration(p.Expires) * time.Second),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})

	if err := uc.codeRepo.StoreAuthorizationCode(code); err != nil {
		return nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	return code, nil
}

func (uc *AuthorizationUsecase) GenerateTokenByCode(code string) (*domain.Token, *domain.RefreshToken, error) {
	var atokn *domain.Token
	var rtokn *domain.RefreshToken
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
	if currentTime.After(c.GetExpiresAt()) {
		return atokn, rtokn, errors.NewUsecaseError(http.StatusForbidden, "code has expired")
	}

	// create token and refresh token
	expiration := time.Now().Add(time.Duration(uc.config.AuthTokenExpiresMin) * time.Minute)
	t := accesstoken.TokenParams{
		UserID:    c.GetUserID(),
		ClientID:  c.GetClientID(),
		Scope:     c.GetScope(),
		ExpiresAt: expiration,
	}
	accessToken, err := uc.tokenGen.Generate(&t, uc.config.PrivateKey)
	if err != nil {
		return atokn, rtokn, errors.NewUsecaseError(http.StatusInternalServerError, "code has expired")
	}

	if err = uc.tokenRepo.StoreToken(&domain.Token{
		AccessToken: accessToken,
		ClientID:    c.GetClientID(),
		UserID:      c.GetUserID(),
		Scope:       c.GetScope(),
		ExpiresAt:   expiration,
	}); err != nil {
		return atokn, rtokn, errors.NewUsecaseError(http.StatusInternalServerError, "code has expired")
	}

	randomString, err := str.GenerateRandomString(randomStringLen)
	if err != nil {
		return atokn, rtokn, errors.NewUsecaseError(http.StatusInternalServerError, "code has expired")
	}
	refreshExpiration := time.Now().Add(time.Duration(uc.config.AuthRefreshTokenExpiresDay) * day)
	if err = uc.refreshTokenRepo.StoreRefreshToken(&domain.RefreshToken{
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

	atokn = &domain.Token{
		AccessToken: accessToken,
		ExpiresAt:   expiration,
	}

	rtokn = &domain.RefreshToken{
		RefreshToken: randomString,
		ExpiresAt:    refreshExpiration,
	}

	return atokn, rtokn, nil
}

func (uc *AuthorizationUsecase) GenerateTokenByRefreshToken(refreshToken string) (*domain.Token, *domain.RefreshToken, error) {
	var atokn *domain.Token
	var rtokn *domain.RefreshToken
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

	if err = uc.tokenRepo.StoreToken(&domain.Token{
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

	if err = uc.refreshTokenRepo.StoreRefreshToken(&domain.RefreshToken{
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

	atokn = &domain.Token{
		AccessToken: accessToken,
		ExpiresAt:   expiration,
	}
	rtokn = &domain.RefreshToken{
		RefreshToken: randomString,
		ExpiresAt:    refreshExpiration,
	}

	return atokn, rtokn, nil
}
