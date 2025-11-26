package domainservice

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/domain"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

//go:generate go run github.com/matryer/moq -out token_service_mock.go . TokenService
type TokenService interface {
	StoreNewToken(ctx context.Context, clientID, UserID uuid.UUID, scope string) (domain.Token, error)
	StoreNewRefreshToken(ctx context.Context, accessToken string) (domain.RefreshToken, error)
	FindToken(ctx context.Context, accessToken string) (domain.Token, error)
	RevokeToken(ctx context.Context, accessToken string) error
	FindRefreshToken(ctx context.Context, refreshToken string) (domain.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, refreshToken string) error
	FindTokenByRefreshToken(ctx context.Context, refreshToken string, now time.Time) (domain.Token, error)
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

func (s *tokenService) StoreNewToken(ctx context.Context, clientID, UserID uuid.UUID, scope string) (domain.Token, error) {
	atoken := domain.NewToken(domain.TokenParams{
		ClientID: clientID,
		UserID:   UserID,
		Scope:    scope,
	})

	var at domain.AccessToken
	token, err := at.Generate(atoken, s.config.PrivateKey)
	if err != nil {
		return nil, err
	}

	if err := atoken.SetNewAccessToken(token); err != nil {
		return nil, err
	}

	atoken.SetNewExpiry(s.config.AuthTokenExpiresMin)

	if err := s.tokenRepo.StoreToken(ctx, atoken); err != nil {
		return nil, err
	}

	return atoken, nil
}

func (s *tokenService) StoreNewRefreshToken(ctx context.Context, accessToken string) (domain.RefreshToken, error) {
	var rt domain.RefreshTokenString
	refreshToken, err := rt.Generate()
	if err != nil {
		return nil, err
	}

	rtoken := domain.NewRefreshToken(domain.RefreshTokenParams{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})

	rtoken.SetNewExpiry(s.config.AuthRefreshTokenExpiresDay)

	if err := s.refreshTokenRepo.StoreRefreshToken(ctx, rtoken); err != nil {
		return nil, err
	}
	return rtoken, nil
}

func (s *tokenService) FindToken(ctx context.Context, accessToken string) (domain.Token, error) {
	return s.tokenRepo.FindToken(ctx, accessToken)
}

func (s *tokenService) RevokeToken(ctx context.Context, accessToken string) error {
	return s.tokenRepo.RevokeToken(ctx, accessToken)
}

func (s *tokenService) FindRefreshToken(ctx context.Context, refreshToken string) (domain.RefreshToken, error) {
	return s.refreshTokenRepo.FindRefreshToken(ctx, refreshToken)
}

func (s *tokenService) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	return s.refreshTokenRepo.RevokeRefreshToken(ctx, refreshToken)
}

func (s *tokenService) FindTokenByRefreshToken(ctx context.Context, refreshToken string, now time.Time) (domain.Token, error) {
	rt, err := s.refreshTokenRepo.FindRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.NewServiceErrorError(errors.ErrCodeInternalServer, err.Error())
	}
	if rt == nil {
		return nil, errors.NewServiceErrorError(errors.ErrCodeNotFound, "refresh token not found")
	}
	if rt.IsExpired(now) {
		return nil, errors.NewServiceErrorError(errors.ErrCodeForbidden, "refresh token has expired")
	}

	tkn, err := s.FindToken(ctx, rt.GetAccessToken())
	if tkn == nil {
		return nil, errors.NewServiceErrorError(errors.ErrCodeNotFound, "token not found")
	}

	return tkn, err
}
