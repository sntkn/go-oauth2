package usecases

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/internal/accesstoken"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestGetUserInvoke(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	t.Run("successful get user", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

		token, err := accesstoken.Generate(accesstoken.TokenParams{
			UserID:    uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			ClientID:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			Scope:     "read",
			ExpiresAt: time.Now().Add(1 * time.Hour),
		})
		assert.NoError(t, err)

		c.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

		mockRepo := &repository.OAuth2RepositoryMock{}
		getUser := NewGetUser(mockRepo)
		mockRepo.FindUserFunc = func(id uuid.UUID) (repository.User, error) {
			return repository.User{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000000"),
				Name:      "client Name",
				Email:     "test@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		}

		user, err := getUser.Invoke(c)

		assert.NoError(t, err)
		assert.Equal(t, "test@example.com", user.Email)
	})

	t.Run("missing authorization header", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

		mockRepo := &repository.OAuth2RepositoryMock{}
		getUser := NewGetUser(mockRepo)
		user, err := getUser.Invoke(c)

		assert.Error(t, err)
		assert.IsType(t, &cerrs.UsecaseError{}, err)
		assert.Equal(t, repository.User{}, user)
		assert.Equal(t, http.StatusUnauthorized, err.(*cerrs.UsecaseError).Code)
		assert.Equal(t, "Code: 401, Message: missing or empty authorization header", err.Error())
	})

	t.Run("missing token", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

		token, err := accesstoken.Generate(accesstoken.TokenParams{
			UserID:    uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			ClientID:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			Scope:     "read",
			ExpiresAt: time.Now().Add(-1 * time.Hour),
		})
		assert.NoError(t, err)

		c.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

		mockRepo := &repository.OAuth2RepositoryMock{}
		getUser := NewGetUser(mockRepo)
		mockRepo.FindUserFunc = func(id uuid.UUID) (repository.User, error) {
			return repository.User{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000000"),
				Name:      "client Name",
				Email:     "test@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		}

		user, err := getUser.Invoke(c)

		assert.Error(t, err)
		assert.IsType(t, &cerrs.UsecaseError{}, err)
		assert.Equal(t, repository.User{}, user)
		assert.Equal(t, http.StatusUnauthorized, err.(*cerrs.UsecaseError).Code)
		assert.Contains(t, err.Error(), "token is expired by")
	})

	t.Run("missing find user", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

		token, err := accesstoken.Generate(accesstoken.TokenParams{
			UserID:    uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			ClientID:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			Scope:     "read",
			ExpiresAt: time.Now().Add(1 * time.Hour),
		})
		assert.NoError(t, err)

		c.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

		mockRepo := &repository.OAuth2RepositoryMock{}
		getUser := NewGetUser(mockRepo)
		mockRepo.FindUserFunc = func(id uuid.UUID) (repository.User, error) {
			return repository.User{}, sql.ErrNoRows
		}

		user, err := getUser.Invoke(c)

		assert.Error(t, err)
		assert.IsType(t, &cerrs.UsecaseError{}, err)
		assert.Equal(t, repository.User{}, user)
		assert.Equal(t, http.StatusUnauthorized, err.(*cerrs.UsecaseError).Code)
		assert.Contains(t, err.Error(), "no rows in result set")
	})
}
