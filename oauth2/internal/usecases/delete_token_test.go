package usecases

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteToken_Invoke(t *testing.T) {
	t.Parallel()
	mockRepo := &repository.OAuth2RepositoryMock{}
	deleteToken := NewDeleteToken(mockRepo)

	gin.SetMode(gin.TestMode)

	t.Run("missing authorization header", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/", http.NoBody)

		err := deleteToken.Invoke(c)

		require.Error(t, err)
		assert.IsType(t, &cerrs.UsecaseError{}, err)
		assert.Equal(t, http.StatusUnauthorized, err.(*cerrs.UsecaseError).Code)
	})

	t.Run("successful token revocation", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/", http.NoBody)
		c.Request.Header.Set("Authorization", "Bearer valid-token")

		mockRepo.RevokeTokenFunc = func(accessToken string) error {
			return nil
		}

		err := deleteToken.Invoke(c)

		require.NoError(t, err)
	})

	t.Run("token revocation error", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/", http.NoBody)
		c.Request.Header.Set("Authorization", "Bearer invalid-token")

		mockRepo.RevokeTokenFunc = func(accessToken string) error {
			return cerrs.NewUsecaseError(http.StatusInternalServerError, "revocation error")
		}

		err := deleteToken.Invoke(c)

		require.Error(t, err)
		assert.IsType(t, &cerrs.UsecaseError{}, err)
		assert.Equal(t, http.StatusInternalServerError, err.(*cerrs.UsecaseError).Code)
	})
}
