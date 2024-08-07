package usecases

import (
	"net/http/httptest"
	"testing"

	cdberr "github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/entity"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestInvoke(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful invoke", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		input := AuthorizeInput{ClientID: "1234-abcd-qwer-asdf"}
		form := entity.SessionSigninForm{Email: "test@example.com", Error: ""}

		mockSess := &session.SessionClientMock{
			GetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				*t.(*AuthorizeInput) = input
				return nil
			},
			FlushNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				*t.(*entity.SessionSigninForm) = form
				return nil
			},
		}
		cfg := &config.Config{}
		signin := NewSignin(cfg, mockSess)

		result, err := signin.Invoke(c)
		assert.NoError(t, err)
		assert.Equal(t, form, result)
	})

	t.Run("missing client_id", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		input := AuthorizeInput{ClientID: ""}
		form := entity.SessionSigninForm{Email: "test@example.com", Error: ""}

		mockSess := &session.SessionClientMock{
			GetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				*t.(*AuthorizeInput) = input
				return nil
			},
			FlushNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				*t.(*entity.SessionSigninForm) = form
				return nil
			},
		}
		cfg := &config.Config{}
		signin := NewSignin(cfg, mockSess)

		result, err := signin.Invoke(c)
		assert.Error(t, err)
		assert.IsType(t, &cerrs.UsecaseError{}, err)
		assert.Contains(t, err.Error(), "invalid client_id")
		assert.Equal(t, entity.SessionSigninForm{}, result)
	})

	t.Run("session data error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		mockSess := &session.SessionClientMock{
			GetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				return cdberr.New("session error")
			},
			FlushNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				return nil
			},
		}
		cfg := &config.Config{}
		signin := NewSignin(cfg, mockSess)

		result, err := signin.Invoke(c)
		assert.Error(t, err)
		assert.IsType(t, &cerrs.UsecaseError{}, err)
		assert.Contains(t, err.Error(), "session error")
		assert.Equal(t, entity.SessionSigninForm{}, result)
	})

	t.Run("flush session data error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		input := AuthorizeInput{ClientID: "1234-abcd-qwer-asdf"}

		mockSess := &session.SessionClientMock{
			GetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				*t.(*AuthorizeInput) = input
				return nil
			},
			FlushNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				return cdberr.New("flush error")
			},
		}
		cfg := &config.Config{}
		signin := NewSignin(cfg, mockSess)

		result, err := signin.Invoke(c)
		assert.Error(t, err)
		assert.IsType(t, &cerrs.UsecaseError{}, err)
		assert.Contains(t, err.Error(), "flush error")
		assert.Equal(t, entity.SessionSigninForm{}, result)
	})
}
