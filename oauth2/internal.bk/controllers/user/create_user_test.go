package user

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/common/session"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUserHandler(t *testing.T) {
	t.Parallel()
	// Ginのテストモードをセット
	gin.SetMode(gin.TestMode)

	t.Run("successful", func(t *testing.T) {
		t.Parallel()
		r := gin.Default()
		r.LoadHTMLGlob("../../../templates/*")

		handler := &CreateUserHandler{
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
							return nil
						},
					}
				},
			},
			uc: &CreateUserUsecaseMock{
				InvokeFunc: func(user *repository.User) error {
					return nil
				},
			},
		}

		r.POST("/signup", handler.CreateUser)

		values := url.Values{}
		values.Add("name", "test")
		values.Add("email", "test@example.com")
		values.Add("password", "test1234")

		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/signup", strings.NewReader(values.Encode()))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusFound, w.Code)
	})

	t.Run("invalid request", func(t *testing.T) {
		t.Parallel()
		r := gin.Default()
		r.LoadHTMLGlob("../../../templates/*")

		handler := &CreateUserHandler{
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
							return nil
						},
					}
				},
			},
			uc: &CreateUserUsecaseMock{},
		}

		r.POST("/signup", handler.CreateUser)

		values := url.Values{}
		values.Add("name", "test")
		values.Add("email", "test@example.com")

		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/signup", strings.NewReader(values.Encode()))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusFound, w.Code)
	})

	t.Run("usecase badrequest error", func(t *testing.T) {
		t.Parallel()
		r := gin.Default()
		r.LoadHTMLGlob("../../../templates/*")

		handler := &CreateUserHandler{
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
							return nil
						},
					}
				},
			},
			uc: &CreateUserUsecaseMock{
				InvokeFunc: func(user *repository.User) error {
					return errors.NewUsecaseError(http.StatusBadRequest, "bad request error")
				},
			},
		}

		r.POST("/signup", handler.CreateUser)

		values := url.Values{}
		values.Add("name", "test")
		values.Add("email", "test@example.com")
		values.Add("password", "test1234")

		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/signup", strings.NewReader(values.Encode()))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusFound, w.Code)
	})

	t.Run("usecase internal server error", func(t *testing.T) {
		t.Parallel()
		r := gin.Default()
		r.LoadHTMLGlob("../../../templates/*")

		handler := &CreateUserHandler{
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
							return nil
						},
					}
				},
			},
			uc: &CreateUserUsecaseMock{
				InvokeFunc: func(user *repository.User) error {
					return errors.NewUsecaseError(http.StatusInternalServerError, "internal server error")
				},
			},
		}

		r.POST("/signup", handler.CreateUser)

		values := url.Values{}
		values.Add("name", "test")
		values.Add("email", "test@example.com")
		values.Add("password", "test1234")

		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/signup", strings.NewReader(values.Encode()))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("internal server error", func(t *testing.T) {
		t.Parallel()
		r := gin.Default()
		r.LoadHTMLGlob("../../../templates/*")

		handler := &CreateUserHandler{
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
							return nil
						},
					}
				},
			},
			uc: &CreateUserUsecaseMock{
				InvokeFunc: func(user *repository.User) error {
					return errors.New("internal server error")
				},
			},
		}

		r.POST("/signup", handler.CreateUser)

		values := url.Values{}
		values.Add("name", "test")
		values.Add("email", "test@example.com")
		values.Add("password", "test1234")

		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/signup", strings.NewReader(values.Encode()))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
