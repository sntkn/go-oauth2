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
	c, err := uc.codeRepo.FindValidAuthorizationCode(code, time.Now())
	if err != nil {
		return nil, nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	if c == nil {
		return nil, nil, errors.NewUsecaseError(http.StatusForbidden, "code not found")
	}

	if c.IsExpired(time.Now()) {
		return nil, nil, errors.NewUsecaseError(http.StatusForbidden, "code has expired")
	}

	atoken, err := uc.storeNewToken(
		c.GetClientID(),
		c.GetUserID(),
		c.GetScope(),
	)
	if err != nil {
		return nil, nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	rtoken, err := uc.storeNewRefreshToken(atoken.GetAccessToken())
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

	atoken, err := uc.storeNewToken(
		tkn.GetClientID(),
		tkn.GetUserID(),
		tkn.GetScope(),
	)
	if err != nil {
		return nil, nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	rtoken, err := uc.storeNewRefreshToken(atoken.GetAccessToken())
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

func (uc *AuthorizationUsecase) storeNewToken(clientID, UserID uuid.UUID, scope string) (domain.Token, error) {
	atoken := domain.NewToken(domain.TokenParams{
		ClientID: clientID,
		UserID:   UserID,
		Scope:    scope,
	})

	if err := atoken.SetNewAccessToken(uc.config.PrivateKey); err != nil {
		return nil, err
	}

	atoken.SetNewExpiry(uc.config.AuthTokenExpiresMin)

	if err := uc.tokenRepo.StoreToken(atoken); err != nil {
		return nil, err
	}

	return atoken, nil
}

func (uc *AuthorizationUsecase) storeNewRefreshToken(accessToken string) (domain.RefreshToken, error) {
	rtoken := domain.NewRefreshToken(domain.RefreshTokenParams{
		AccessToken: accessToken,
	})
	if err := rtoken.SetNewRefreshToken(); err != nil {
		return nil, err
	}
	rtoken.SetNewExpiry(uc.config.AuthRefreshTokenExpiresDay)

	if err := uc.refreshTokenRepo.StoreRefreshToken(rtoken); err != nil {
		return nil, err
	}
	return rtoken, nil
}
