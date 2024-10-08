package user

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/entity"
	"github.com/sntkn/go-oauth2/oauth2/internal/flashmessage"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockSignupFinishedUsecase struct {
	mock.Mock
}

func (m *MockSignupFinishedUsecase) Invoke(c *gin.Context) (entity.SessionRegistrationForm, error) {
	args := m.Called(c)
	return args.Get(0).(entity.SessionRegistrationForm), args.Error(1)
}

func TestSignupFinishedHandler(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	t.Run("signup finished successful", func(t *testing.T) {
		t.Parallel()
		r := gin.Default()
		r.LoadHTMLGlob("../../../templates/*")

		r.Use(func(c *gin.Context) {
			c.Set("flashMessages", &flashmessage.Messages{})
			c.Next()
		})

		handler := &SignupFinishedHandler{
			sessionManager: &session.SessionManagerMock{
				NewSessionFunc: func(c *gin.Context) session.SessionClient {
					return &session.SessionClientMock{
						SetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
							return nil
						},
						DelSessionDataFunc: func(c *gin.Context, key string) error {
							return nil
						},
						GetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
							*t.(*flashmessage.Messages) = flashmessage.Messages{}
							return nil
						},
					}
				},
			},
		}

		r.GET("/signup_finished", handler.SignupFinished)

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/signup_finished", http.NoBody)
		require.NoError(t, err)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		assert.Contains(t, w.Body.String(), "User creation was successful.")
	})
}
