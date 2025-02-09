package usecase

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/domain"
	"github.com/sntkn/go-oauth2/oauth2/domain/domainservice"
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
	tokenService domainservice.TokenService,
) IAuthorizationUsecase {
	return &AuthorizationUsecase{
		clientRepo:   clientRepo,
		codeRepo:     codeRepo,
		tokenService: tokenService,
	}
}

type AuthorizationUsecase struct {
	clientRepo   domain.ClientRepository
	codeRepo     domain.AuthorizationCodeRepository
	tokenService domainservice.TokenService
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

	code, err := uc.codeRepo.StoreAuthorizationCode(domain.StoreAuthorizationCodeParams{
		Code:        randomString,
		Scope:       p.Scope,
		RedirectURI: p.RedirectURI,
		ExpiresAt:   time.Now().Add(time.Duration(p.Expires) * time.Second),
	})
	if err != nil {
		return nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	c, err := uc.codeRepo.FindAuthorizationCode(code)
	if err != nil {
		return nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}
	if c == nil {
		return nil, errors.NewUsecaseError(http.StatusBadRequest, "code not found")
	}

	return c, nil
}

func (uc *AuthorizationUsecase) GenerateTokenByCode(code string) (domain.Token, domain.RefreshToken, error) {
	c, err := uc.codeRepo.FindAuthorizationCode(code)
	if err != nil {
		return nil, nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	if c == nil {
		return nil, nil, errors.NewUsecaseError(http.StatusForbidden, "code not found")
	}

	if c.IsExpired(time.Now()) {
		return nil, nil, errors.NewUsecaseError(http.StatusForbidden, "code has expired")
	}

	atoken, err := uc.tokenService.StoreNewToken(
		c.GetClientID(),
		c.GetUserID(),
		c.GetScope(),
	)
	if err != nil {
		return nil, nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	rtoken, err := uc.tokenService.StoreNewRefreshToken(atoken.GetAccessToken())
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

	tkn, rt, err := uc.tokenService.FindTokenAndRefreshTokenByRefreshToken(refreshToken, time.Now())
	if err != nil {
		if serviceErr, ok := err.(*errors.ServiceError); ok {
			return nil, nil, errors.NewUsecaseError(serviceErr.Code, serviceErr.Error())
		}
		return nil, nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	atoken, err := uc.tokenService.StoreNewToken(
		tkn.GetClientID(),
		tkn.GetUserID(),
		tkn.GetScope(),
	)
	if err != nil {
		return nil, nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	rtoken, err := uc.tokenService.StoreNewRefreshToken(atoken.GetAccessToken())
	if err != nil {
		return nil, nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	if err := uc.tokenService.RevokeToken(tkn.GetAccessToken()); err != nil {
		return nil, nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	if err := uc.tokenService.RevokeRefreshToken(rt.GetRefreshToken()); err != nil {
		return nil, nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	return atoken, rtoken, nil
}
