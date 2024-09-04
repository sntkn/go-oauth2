package usecases

import (
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTokenByCode_Invoke(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	t.Run("valid token", func(t *testing.T) {
		t.Parallel()
		mockRepo := &repository.OAuth2RepositoryMock{
			FindValidOAuth2CodeFunc: func(code string, expiresAt time.Time) (repository.Code, error) {
				return repository.Code{
					UserID:    uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					ClientID:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					Scope:     "",
					ExpiresAt: time.Now().Add(1 * time.Hour),
				}, nil
			},
			RegisterTokenFunc: func(t *repository.Token) error {
				return nil
			},
			RegisterRefreshTokenFunc: func(t *repository.RefreshToken) error {
				return nil
			},
			RevokeCodeFunc: func(code string) error {
				return nil
			},
		}
		u := &CreateTokenByCode{
			db:  mockRepo,
			cfg: &config.Config{},
		}

		authCode := "valid_auth_code"
		c, _ := gin.CreateTestContext(nil)
		authTokens, err := u.Invoke(c, authCode)

		require.NoError(t, err)
		assert.NotEmpty(t, authTokens.AccessToken)
		assert.NotEmpty(t, authTokens.RefreshToken)
		assert.NotZero(t, authTokens.Expiry)
	})

	t.Run("invalid auth code", func(t *testing.T) {
		t.Parallel()
		mockRepo := &repository.OAuth2RepositoryMock{
			FindValidOAuth2CodeFunc: func(code string, expiresAt time.Time) (repository.Code, error) {
				return repository.Code{}, sql.ErrNoRows
			},
		}
		u := &CreateTokenByCode{
			db:  mockRepo,
			cfg: &config.Config{},
		}

		authCode := "invalid_auth_code"

		c, _ := gin.CreateTestContext(nil)
		authTokens, err := u.Invoke(c, authCode)

		require.Error(t, err)
		assert.Equal(t, http.StatusForbidden, err.(*errors.UsecaseError).Code)
		assert.Empty(t, authTokens.AccessToken)
		assert.Empty(t, authTokens.RefreshToken)
	})

	t.Run("database error on finding refresh token", func(t *testing.T) {
		t.Parallel()
		mockRepo := &repository.OAuth2RepositoryMock{
			FindValidOAuth2CodeFunc: func(code string, expiresAt time.Time) (repository.Code, error) {
				return repository.Code{}, errors.New("db error")
			},
		}
		u := &CreateTokenByCode{
			db:  mockRepo,
			cfg: &config.Config{},
		}
		authCode := "db_error_auth_code"

		c, _ := gin.CreateTestContext(nil)
		authTokens, err := u.Invoke(c, authCode)

		require.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.(*errors.UsecaseError).Code)
		assert.Empty(t, authTokens.AccessToken)
		assert.Empty(t, authTokens.RefreshToken)
	})
}
