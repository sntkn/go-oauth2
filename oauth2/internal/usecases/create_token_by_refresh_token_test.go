package usecases

import (
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestCreateTokenByRefreshToken_Invoke(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("valid refresh token", func(t *testing.T) {
		mockRepo := &repository.OAuth2RepositoryMock{
			FindValidRefreshTokenFunc: func(refreshToken string, expiresAt time.Time) (repository.RefreshToken, error) {
				return repository.RefreshToken{
					AccessToken:  "valid_access_token",
					RefreshToken: "valid_refresh_token",
					ExpiresAt:    time.Now().Add(1 * time.Hour),
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				}, nil
			},
			FindTokenFunc: func(accessToken string) (repository.Token, error) {
				return repository.Token{
					AccessToken: "valid_access_token",
					UserID:      uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					ClientID:    uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					Scope:       "scope",
					ExpiresAt:   time.Now().Add(1 * time.Hour),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}, nil
			},
			RegisterTokenFunc: func(t *repository.Token) error {
				return nil
			},
			RegisterRefreshTokenFunc: func(t *repository.RefreshToken) error {
				return nil
			},
			RevokeTokenFunc: func(accessToken string) error {
				return nil
			},
			RevokeRefreshTokenFunc: func(refreshToken string) error {
				return nil
			},
		}
		u := &CreateTokenByRefreshToken{
			db:  mockRepo,
			cfg: &config.Config{},
		}

		refreshToken := "valid_refresh_token"
		c, _ := gin.CreateTestContext(nil)
		authTokens, err := u.Invoke(c, refreshToken)

		assert.NoError(t, err)
		assert.NotEmpty(t, authTokens.AccessToken)
		assert.NotEmpty(t, authTokens.RefreshToken)
		assert.NotZero(t, authTokens.Expiry)
	})

	t.Run("invalid refresh token", func(t *testing.T) {
		mockRepo := &repository.OAuth2RepositoryMock{
			FindValidRefreshTokenFunc: func(refreshToken string, expiresAt time.Time) (repository.RefreshToken, error) {
				return repository.RefreshToken{}, sql.ErrNoRows
			},
		}
		u := &CreateTokenByRefreshToken{
			db:  mockRepo,
			cfg: &config.Config{},
		}

		refreshToken := "invalid_refresh_token"

		c, _ := gin.CreateTestContext(nil)
		authTokens, err := u.Invoke(c, refreshToken)

		assert.Error(t, err)
		assert.Equal(t, http.StatusForbidden, err.(*cerrs.UsecaseError).Code)
		assert.Empty(t, authTokens.AccessToken)
		assert.Empty(t, authTokens.RefreshToken)
	})

	t.Run("database error on finding refresh token", func(t *testing.T) {
		mockRepo := &repository.OAuth2RepositoryMock{
			FindValidRefreshTokenFunc: func(refreshToken string, expiresAt time.Time) (repository.RefreshToken, error) {
				return repository.RefreshToken{}, errors.New("db error")
			},
		}
		u := &CreateTokenByRefreshToken{
			db:  mockRepo,
			cfg: &config.Config{},
		}
		refreshToken := "db_error_refresh_token"

		c, _ := gin.CreateTestContext(nil)
		authTokens, err := u.Invoke(c, refreshToken)

		assert.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.(*cerrs.UsecaseError).Code)
		assert.Empty(t, authTokens.AccessToken)
		assert.Empty(t, authTokens.RefreshToken)
	})

	t.Run("invalid access token", func(t *testing.T) {
		mockRepo := &repository.OAuth2RepositoryMock{
			FindValidRefreshTokenFunc: func(refreshToken string, expiresAt time.Time) (repository.RefreshToken, error) {
				return repository.RefreshToken{
					AccessToken:  "valid_access_token",
					RefreshToken: "valid_refresh_token",
					ExpiresAt:    time.Now().Add(1 * time.Hour),
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				}, nil
			},
			FindTokenFunc: func(accessToken string) (repository.Token, error) {
				return repository.Token{}, sql.ErrNoRows
			},
		}
		u := &CreateTokenByRefreshToken{
			db:  mockRepo,
			cfg: &config.Config{},
		}
		refreshToken := "valid_refresh_token"

		c, _ := gin.CreateTestContext(nil)
		authTokens, err := u.Invoke(c, refreshToken)

		assert.Error(t, err)
		assert.Equal(t, http.StatusForbidden, err.(*cerrs.UsecaseError).Code)
		assert.Empty(t, authTokens.AccessToken)
		assert.Empty(t, authTokens.RefreshToken)
	})

	t.Run("database error on finding access token", func(t *testing.T) {
		mockRepo := &repository.OAuth2RepositoryMock{
			FindValidRefreshTokenFunc: func(refreshToken string, expiresAt time.Time) (repository.RefreshToken, error) {
				return repository.RefreshToken{
					AccessToken:  "valid_access_token",
					RefreshToken: "valid_refresh_token",
					ExpiresAt:    time.Now().Add(1 * time.Hour),
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				}, nil
			},
			FindTokenFunc: func(accessToken string) (repository.Token, error) {
				return repository.Token{}, errors.New("db error")
			},
		}
		u := &CreateTokenByRefreshToken{
			db:  mockRepo,
			cfg: &config.Config{},
		}
		refreshToken := "valid_refresh_token"

		c, _ := gin.CreateTestContext(nil)
		authTokens, err := u.Invoke(c, refreshToken)

		assert.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.(*cerrs.UsecaseError).Code)
		assert.Empty(t, authTokens.AccessToken)
		assert.Empty(t, authTokens.RefreshToken)
	})
}
