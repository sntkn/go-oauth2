package usecases

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteToken_Invoke(t *testing.T) {
	t.Parallel()
	mockRepo := &repository.OAuth2RepositoryMock{}
	deleteToken := NewDeleteToken(mockRepo)

	gin.SetMode(gin.TestMode)

	t.Run("successful token revocation", func(t *testing.T) {
		t.Parallel()

		mockRepo.RevokeTokenFunc = func(accessToken string) error {
			return nil
		}

		err := deleteToken.Invoke("token")

		require.NoError(t, err)
	})

	t.Run("token revocation error", func(t *testing.T) {
		t.Parallel()

		mockRepo.RevokeTokenFunc = func(accessToken string) error {
			return errors.NewUsecaseError(http.StatusInternalServerError, "revocation error")
		}

		err := deleteToken.Invoke("token")

		require.Error(t, err)
		assert.IsType(t, &errors.UsecaseError{}, err)
		assert.Equal(t, http.StatusInternalServerError, err.(*errors.UsecaseError).Code)
	})
}
