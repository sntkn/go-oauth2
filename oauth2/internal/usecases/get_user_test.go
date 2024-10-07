package usecases

import (
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserInvoke(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	t.Run("successful get user", func(t *testing.T) {
		t.Parallel()

		userID := uuid.MustParse("00000000-0000-0000-0000-000000000000")

		mockRepo := &repository.OAuth2RepositoryMock{}
		getUser := NewGetUser(mockRepo)
		mockRepo.FindUserFunc = func(id uuid.UUID) (repository.User, error) {
			return repository.User{
				ID:        userID,
				Name:      "client Name",
				Email:     "test@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		}

		user, err := getUser.Invoke(userID)

		require.NoError(t, err)
		assert.Equal(t, "test@example.com", user.Email)
	})

	t.Run("missing find user", func(t *testing.T) {
		t.Parallel()

		userID := uuid.MustParse("00000000-0000-0000-0000-000000000000")

		mockRepo := &repository.OAuth2RepositoryMock{}
		getUser := NewGetUser(mockRepo)
		mockRepo.FindUserFunc = func(id uuid.UUID) (repository.User, error) {
			return repository.User{}, sql.ErrNoRows
		}

		user, err := getUser.Invoke(userID)

		require.Error(t, err)
		assert.IsType(t, &errors.UsecaseError{}, err)
		assert.Equal(t, repository.User{}, user)
		assert.Equal(t, http.StatusUnauthorized, err.(*errors.UsecaseError).Code)
		assert.Contains(t, err.Error(), "no rows in result set")
	})
}
