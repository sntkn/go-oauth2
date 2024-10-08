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
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockSignupUsecase struct {
	mock.Mock
}

func (m *MockSignupUsecase) Invoke(c *gin.Context) (entity.SessionRegistrationForm, error) {
	args := m.Called(c)
	return args.Get(0).(entity.SessionRegistrationForm), args.Error(1)
}

func TestSignupHandler(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.LoadHTMLGlob("../../../templates/*")

	form := entity.SessionRegistrationForm{
		Name:  "test",
		Email: "test@example.com",
		Error: "",
	}
	handler := &SignupHandler{
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
						switch v := t.(type) {
						case *entity.SessionRegistrationForm:
							*v = form
						case *flashmessage.Messages:
							*v = flashmessage.Messages{}
						default:
							return errors.New("interface conversion error")
						}
						return nil
					},
					FlushNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
						return nil
					},
				}
			},
		},
	}

	r.GET("/signup", handler.Signup)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/signup", http.NoBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	assert.Contains(t, w.Body.String(), "<h2>Signup</h2>")
}
