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
	"golang.org/x/crypto/bcrypt"
)

func TestAuthorizationInvoke(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	t.Run("successful invoke", func(t *testing.T) {
		t.Parallel()

		mockRepo := &repository.OAuth2RepositoryMock{
			FindUserByEmailFunc: func(email string) (repository.User, error) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("test1234"), bcrypt.DefaultCost)
				return repository.User{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					Name:      "client Name",
					Email:     "test@example.com",
					Password:  string(hashedPassword),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
			RegisterOAuth2CodeFunc: func(c *repository.Code) error {
				return nil
			},
		}
		cfg := &config.Config{}
		authorization := NewAuthorization(cfg, mockRepo)
		params := AuthorizationInput{
			Email:       "test@example.com",
			Password:    "test1234",
			ClientID:    "00000000-0000-0000-0000-000000000000",
			Scope:       "read",
			RedirectURI: "http://example.com",
			Expires:     0,
		}

		result, err := authorization.Invoke(params)
		require.NoError(t, err)
		assert.Contains(t, result, "http://example.com?code=")
	})

	t.Run("user not found error", func(t *testing.T) {
		t.Parallel()

		mockRepo := &repository.OAuth2RepositoryMock{
			FindUserByEmailFunc: func(email string) (repository.User, error) {
				return repository.User{}, sql.ErrNoRows
			},
			RegisterOAuth2CodeFunc: func(c *repository.Code) error {
				return nil
			},
		}
		cfg := &config.Config{}
		authorization := NewAuthorization(cfg, mockRepo)
		params := AuthorizationInput{
			Email:       "test@example.com",
			Password:    "test1234",
			ClientID:    "00000000-0000-0000-0000-000000000000",
			Scope:       "read",
			RedirectURI: "http://example.com",
			Expires:     0,
		}

		result, err := authorization.Invoke(params)
		require.Error(t, err)
		assert.IsType(t, &errors.UsecaseError{}, err)
		assert.Contains(t, err.Error(), "no rows in result set")
		assert.Equal(t, "", result)
	})

	t.Run("hashed password does not match error", func(t *testing.T) {
		t.Parallel()

		mockRepo := &repository.OAuth2RepositoryMock{
			FindUserByEmailFunc: func(email string) (repository.User, error) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("test1234"), bcrypt.DefaultCost)
				return repository.User{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					Name:      "client Name",
					Email:     "test@example.com",
					Password:  string(hashedPassword),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
			RegisterOAuth2CodeFunc: func(c *repository.Code) error {
				return nil
			},
		}
		cfg := &config.Config{}
		authorization := NewAuthorization(cfg, mockRepo)
		params := AuthorizationInput{
			Email:       "test@example.com",
			Password:    "test12345",
			ClientID:    "00000000-0000-0000-0000-000000000000",
			Scope:       "read",
			RedirectURI: "http://example.com",
			Expires:     0,
		}

		result, err := authorization.Invoke(params)
		require.Error(t, err)
		assert.IsType(t, &errors.UsecaseError{}, err)
		usecaseErr, ok := err.(*errors.UsecaseError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, usecaseErr.Code)
		assert.Equal(t, "", result)
	})

	t.Run("could not parse client_id error", func(t *testing.T) {
		t.Parallel()

		mockRepo := &repository.OAuth2RepositoryMock{
			FindUserByEmailFunc: func(email string) (repository.User, error) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("test1234"), bcrypt.DefaultCost)
				return repository.User{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					Name:      "client Name",
					Email:     "test@example.com",
					Password:  string(hashedPassword),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
			RegisterOAuth2CodeFunc: func(c *repository.Code) error {
				return nil
			},
		}
		cfg := &config.Config{}
		authorization := NewAuthorization(cfg, mockRepo)
		params := AuthorizationInput{
			Email:       "test@example.com",
			Password:    "test1234",
			ClientID:    "00000000-0000-0000-0000-0000000000001",
			Scope:       "read",
			RedirectURI: "http://example.com",
			Expires:     0,
		}

		result, err := authorization.Invoke(params)
		require.Error(t, err)
		assert.IsType(t, &errors.UsecaseError{}, err)
		usecaseErr, ok := err.(*errors.UsecaseError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, usecaseErr.Code)
		assert.Equal(t, "", result)
	})
}
