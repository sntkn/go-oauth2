package usecases

import (
	"net/http/httptest"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUserInvoke(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	t.Run("successful invoke", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		mockRepo := &repository.OAuth2RepositoryMock{
			ExistsUserByEmailFunc: func(email string) (bool, error) {
				return false, nil
			},
			CreateUserFunc: func(u *repository.User) error {
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
		signup := NewCreateUser(cfg, mockRepo, mockSess)

		err := signup.Invoke(c, repository.User{})
		require.NoError(t, err)
	})

	t.Run("email already exists error", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		mockRepo := &repository.OAuth2RepositoryMock{
			ExistsUserByEmailFunc: func(email string) (bool, error) {
				return true, nil
			},
			CreateUserFunc: func(u *repository.User) error {
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
		signup := NewCreateUser(cfg, mockRepo, mockSess)

		err := signup.Invoke(c, repository.User{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "input email already exists")
	})
	t.Run("successful invoke", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		mockRepo := &repository.OAuth2RepositoryMock{
			ExistsUserByEmailFunc: func(email string) (bool, error) {
				return true, errors.New("internal error")
			},
			CreateUserFunc: func(u *repository.User) error {
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
		signup := NewCreateUser(cfg, mockRepo, mockSess)

		err := signup.Invoke(c, repository.User{})
		require.Error(t, err)
		assert.IsType(t, &cerrs.UsecaseError{}, err)
		assert.Contains(t, err.Error(), "internal error")
	})
}
