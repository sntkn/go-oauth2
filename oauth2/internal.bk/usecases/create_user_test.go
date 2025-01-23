package usecases

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUserInvoke(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	t.Run("successful invoke", func(t *testing.T) {
		t.Parallel()

		mockRepo := &repository.OAuth2RepositoryMock{
			ExistsUserByEmailFunc: func(email string) (bool, error) {
				return false, nil
			},
			CreateUserFunc: func(u *repository.User) error {
				return nil
			},
		}
		signup := NewCreateUser(mockRepo)

		err := signup.Invoke(&repository.User{})
		require.NoError(t, err)
	})

	t.Run("email already exists error", func(t *testing.T) {
		t.Parallel()

		mockRepo := &repository.OAuth2RepositoryMock{
			ExistsUserByEmailFunc: func(email string) (bool, error) {
				return true, nil
			},
			CreateUserFunc: func(u *repository.User) error {
				return nil
			},
		}
		signup := NewCreateUser(mockRepo)

		err := signup.Invoke(&repository.User{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "input email already exists")
	})

	t.Run("internal server error invoke", func(t *testing.T) {
		t.Parallel()

		mockRepo := &repository.OAuth2RepositoryMock{
			ExistsUserByEmailFunc: func(email string) (bool, error) {
				return true, errors.New("internal error")
			},
			CreateUserFunc: func(u *repository.User) error {
				return nil
			},
		}
		signup := NewCreateUser(mockRepo)

		err := signup.Invoke(&repository.User{})
		require.Error(t, err)
		assert.IsType(t, &errors.UsecaseError{}, err)
		assert.Contains(t, err.Error(), "internal error")
	})
}
