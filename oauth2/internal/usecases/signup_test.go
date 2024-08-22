package usecases

import (
	"net/http/httptest"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/entity"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestSignupInvoke(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful invoke", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		form := entity.SessionRegistrationForm{
			Name:  "test",
			Email: "test@example.com",
			Error: "",
		}
		mockSess := &session.SessionClientMock{
			FlushNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				*t.(*entity.SessionRegistrationForm) = form
				return nil
			},
		}
		cfg := &config.Config{}
		signup := NewSignup(cfg, mockSess)

		result, err := signup.Invoke(c)
		assert.NoError(t, err)
		assert.Equal(t, form, result)
	})

	t.Run("flush session data error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		form := entity.SessionRegistrationForm{}
		mockSess := &session.SessionClientMock{
			FlushNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				*t.(*entity.SessionRegistrationForm) = form
				return errors.New("flush error")
			},
		}
		cfg := &config.Config{}
		signup := NewSignup(cfg, mockSess)

		result, err := signup.Invoke(c)
		assert.Error(t, err)
		assert.IsType(t, &cerrs.UsecaseError{}, err)
		assert.Contains(t, err.Error(), "flush error")
		assert.Equal(t, entity.SessionRegistrationForm{}, result)
	})
}