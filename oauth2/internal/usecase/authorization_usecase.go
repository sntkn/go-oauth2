package usecase

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/domain"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

type IAuthorizationUsecase interface {
	Consent(uuid.UUID) (domain.Client, error)
	GenerateAuthorizationCode(GenerateAuthorizationCodeParams) (domain.AuthorizationCode, error)
	GenerateTokenByCode(string) (domain.Token, domain.RefreshToken, error)
	GenerateTokenByRefreshToken(string) (domain.Token, domain.RefreshToken, error)
	// GenerateAuthorizationCode(user *model.User, client *model.Client, scopes []string) (*model.AuthorizationCode, error)
	// ValidateAuthorizationCode(code string, clientID string) (*model.AuthorizationCode, error)
}

func NewAuthorizationUsecase(
	clientRepo domain.ClientRepository,
	codeRepo domain.AuthorizationCodeRepository,
	tokenRepo domain.TokenRepository,
	refreshTokenRepo domain.RefreshTokenRepository,
	config *config.Config,
) IAuthorizationUsecase {
	return &AuthorizationUsecase{
		clientRepo:       clientRepo,
		codeRepo:         codeRepo,
		tokenRepo:        tokenRepo,
		refreshTokenRepo: refreshTokenRepo,
		config:           config,
	}
}

type AuthorizationUsecase struct {
	clientRepo       domain.ClientRepository
	codeRepo         domain.AuthorizationCodeRepository
	tokenRepo        domain.TokenRepository
	refreshTokenRepo domain.RefreshTokenRepository
	config           *config.Config
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

func (uc *AuthorizationUsecase) GenerateTokenByCode(code string) (domain.Token, domain.RefreshToken, error) {
	c, err := uc.codeRepo.FindValidAuthorizationCode(code, time.Now())
	if err != nil {
		return nil, nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	if c.IsNotFound() {
		return nil, nil, errors.NewUsecaseError(http.StatusForbidden, "code not found")
	}

	if c.IsExpired(time.Now()) {
		return nil, nil, errors.NewUsecaseError(http.StatusForbidden, "code has expired")
	}

	atoken, err := domain.StoreNewToken(domain.StoreNewTokenParams{
		ClientID:         c.GetClientID(),
		UserID:           c.GetUserID(),
		Scope:            c.GetScope(),
		PrivateKeyBase64: uc.config.PrivateKey,
		AdditionalMin:    uc.config.AuthTokenExpiresMin,
		Repo:             uc.tokenRepo,
	})
	if err != nil {
		return nil, nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	rtoken, err := domain.StoreNewRefreshToken(domain.StoreNewRefreshTokenParams{
		AccessToken:   atoken.GetAccessToken(),
		AdditionalDay: uc.config.AuthRefreshTokenExpiresDay,
		Repo:          uc.refreshTokenRepo,
	})
	if err != nil {
		return nil, nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	// revoke code
	if err := uc.codeRepo.RevokeCode(code); err != nil {
		return nil, nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	return atoken, rtoken, nil
}

func (uc *AuthorizationUsecase) GenerateTokenByRefreshToken(refreshToken string) (domain.Token, domain.RefreshToken, error) {

	rt, err := uc.refreshTokenRepo.FindValidRefreshToken(refreshToken, time.Now())
	if err != nil {
		return nil, nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	if rt.IsNotFound() {
		return nil, nil, errors.NewUsecaseError(http.StatusForbidden, "refresh token not found")
	}

	tkn, err := uc.tokenRepo.FindToken(rt.GetAccessToken())
	if err != nil {
		return nil, nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}
	if tkn.IsNotFound() {
		return nil, nil, errors.NewUsecaseError(http.StatusForbidden, "token not found")
	}

	atoken, err := domain.StoreNewToken(domain.StoreNewTokenParams{
		ClientID:         tkn.GetClientID(),
		UserID:           tkn.GetUserID(),
		Scope:            tkn.GetScope(),
		PrivateKeyBase64: uc.config.PrivateKey,
		AdditionalMin:    uc.config.AuthTokenExpiresMin,
		Repo:             uc.tokenRepo,
	})
	if err != nil {
		return nil, nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	rtoken, err := domain.StoreNewRefreshToken(domain.StoreNewRefreshTokenParams{
		AccessToken:   atoken.GetAccessToken(),
		AdditionalDay: uc.config.AuthRefreshTokenExpiresDay,
		Repo:          uc.refreshTokenRepo,
	})
	if err != nil {
		return nil, nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	if err := uc.tokenRepo.RevokeToken(tkn.GetAccessToken()); err != nil {
		return nil, nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	if err := uc.refreshTokenRepo.RevokeRefreshToken(rt.GetRefreshToken()); err != nil {
		return nil, nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	return atoken, rtoken, nil
}
