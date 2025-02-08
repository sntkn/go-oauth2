package domainservice

import (
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/domain"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
)

//go:generate go run github.com/matryer/moq -out token_service_mock.go . TokenService
type TokenService interface {
	StoreNewToken(clientID, UserID uuid.UUID, scope string) (domain.Token, error)
	StoreNewRefreshToken(accessToken string) (domain.RefreshToken, error)
	FindToken(accessToken string) (domain.Token, error)
	RevokeToken(accessToken string) error
	FindRefreshToken(refreshToken string) (domain.RefreshToken, error)
	RevokeRefreshToken(refreshToken string) error
}

func NewTokenService(
	tokenRepo domain.TokenRepository,
	refreshTokenRepo domain.RefreshTokenRepository,
	config *config.Config,
) *tokenService {
	return &tokenService{
		tokenRepo:        tokenRepo,
		refreshTokenRepo: refreshTokenRepo,
		config:           config,
	}
}

type tokenService struct {
	tokenRepo        domain.TokenRepository
	refreshTokenRepo domain.RefreshTokenRepository
	config           *config.Config
}

func (s *tokenService) StoreNewToken(clientID, UserID uuid.UUID, scope string) (domain.Token, error) {
	atoken := domain.NewToken(domain.TokenParams{
		ClientID: clientID,
		UserID:   UserID,
		Scope:    scope,
	})

	if err := atoken.SetNewAccessToken(s.config.PrivateKey); err != nil {
		return nil, err
	}

	atoken.SetNewExpiry(s.config.AuthTokenExpiresMin)

	if err := s.tokenRepo.StoreToken(atoken); err != nil {
		return nil, err
	}

	return atoken, nil
}

func (s *tokenService) StoreNewRefreshToken(accessToken string) (domain.RefreshToken, error) {
	rtoken := domain.NewRefreshToken(domain.RefreshTokenParams{
		AccessToken: accessToken,
	})
	if err := rtoken.SetNewRefreshToken(); err != nil {
		return nil, err
	}
	rtoken.SetNewExpiry(s.config.AuthRefreshTokenExpiresDay)

	if err := s.refreshTokenRepo.StoreRefreshToken(rtoken); err != nil {
		return nil, err
	}
	return rtoken, nil
}

func (s *tokenService) FindToken(accessToken string) (domain.Token, error) {
	return s.tokenRepo.FindToken(accessToken)
}

func (s *tokenService) RevokeToken(accessToken string) error {
	return s.tokenRepo.RevokeToken(accessToken)
}

func (s *tokenService) FindRefreshToken(refreshToken string) (domain.RefreshToken, error) {
	return s.refreshTokenRepo.FindRefreshToken(refreshToken)
}

func (s *tokenService) RevokeRefreshToken(refreshToken string) error {
	return s.refreshTokenRepo.RevokeRefreshToken(refreshToken)
}
