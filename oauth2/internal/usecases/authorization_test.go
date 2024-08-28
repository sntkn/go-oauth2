package usecases

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthorizationInvoke(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	t.Run("successful invoke", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

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
		mockSess := &session.SessionClientMock{
			DelSessionDataFunc: func(c *gin.Context, key string) error {
				return nil
			},
			SetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				return nil
			},
			GetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				*t.(*AuthorizeInput) = AuthorizeInput{
					ClientID:    "00000000-0000-0000-0000-000000000000",
					Scope:       "read",
					RedirectURI: "http://example.com",
				}
				return nil
			},
		}
		cfg := &config.Config{}
		authorization := NewAuthorization(cfg, mockRepo, mockSess)

		result, err := authorization.Invoke(c, "test@example.com", "test1234")
		assert.NoError(t, err)
		assert.Contains(t, result, "http://example.com?code=")
	})

	t.Run("user not found error", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		mockRepo := &repository.OAuth2RepositoryMock{
			FindUserByEmailFunc: func(email string) (repository.User, error) {
				return repository.User{}, sql.ErrNoRows
			},
			RegisterOAuth2CodeFunc: func(c *repository.Code) error {
				return nil
			},
		}
		mockSess := &session.SessionClientMock{
			DelSessionDataFunc: func(c *gin.Context, key string) error {
				return nil
			},
			SetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				return nil
			},
			GetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				*t.(*AuthorizeInput) = AuthorizeInput{
					ClientID:    "00000000-0000-0000-0000-000000000000",
					Scope:       "read",
					RedirectURI: "http://example.com",
				}
				return nil
			},
		}
		cfg := &config.Config{}
		authorization := NewAuthorization(cfg, mockRepo, mockSess)

		result, err := authorization.Invoke(c, "test@example.com", "test1234")
		assert.Error(t, err)
		assert.IsType(t, &cerrs.UsecaseError{}, err)
		assert.Contains(t, err.Error(), "no rows in result set")
		assert.Equal(t, "", result)
	})

	t.Run("hashed password does not match error", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

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
		mockSess := &session.SessionClientMock{
			DelSessionDataFunc: func(c *gin.Context, key string) error {
				return nil
			},
			SetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				return nil
			},
			GetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				*t.(*AuthorizeInput) = AuthorizeInput{
					ClientID:    "00000000-0000-0000-0000-000000000000",
					Scope:       "read",
					RedirectURI: "http://example.com",
				}
				return nil
			},
		}
		cfg := &config.Config{}
		authorization := NewAuthorization(cfg, mockRepo, mockSess)

		result, err := authorization.Invoke(c, "test@example.com", "test12345")
		assert.Error(t, err)
		assert.IsType(t, &cerrs.UsecaseError{}, err)
		usecaseErr, ok := err.(*cerrs.UsecaseError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, usecaseErr.Code)
		assert.Equal(t, "", result)
	})

	t.Run("could not parse client_id error", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

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
		mockSess := &session.SessionClientMock{
			DelSessionDataFunc: func(c *gin.Context, key string) error {
				return nil
			},
			SetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				return nil
			},
			GetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				return nil
			},
		}
		cfg := &config.Config{}
		authorization := NewAuthorization(cfg, mockRepo, mockSess)

		result, err := authorization.Invoke(c, "test@example.com", "test1234")
		assert.Error(t, err)
		assert.IsType(t, &cerrs.UsecaseError{}, err)
		usecaseErr, ok := err.(*cerrs.UsecaseError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, usecaseErr.Code)
		assert.Equal(t, "", result)
	})
}
