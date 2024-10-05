package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthorize(t *testing.T) {
	t.Parallel()
	// Ginのテストモードをセット
	gin.SetMode(gin.TestMode)

	t.Run("successful authorize", func(t *testing.T) {
		t.Parallel()

		r := gin.Default()
		r.LoadHTMLGlob("../../../templates/*")

		handler := &AuthorizeHandler{
			sessionManager: &session.SessionManagerMock{
				NewSessionFunc: func(c *gin.Context) session.SessionClient {
					return &session.SessionClientMock{
						SetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
							return nil
						},
					}
				},
			},
			uc: &AuthorizeUsecaseMock{
				InvokeFunc: func(clientID string, redirectURI string) error {
					return nil
				},
			},
		}

		r.GET("/authorize", handler.Authorize)

		url := "/authorize?response_type=code&client_id=00000000-0000-0000-0000-000000000000&scope=read&redirect_uri=http://example.com&state=xyz"

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
		require.NoError(t, err)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusFound, w.Code)
	})

	invalidParams := []string{
		"client_id=00000000-0000-0000-0000-000000000000&scope=read&redirect_uri=http://example.com&state=xyz",          // response_typeが欠けている
		"response_type=code&scope=read&redirect_uri=http://example.com&state=xyz",                                      // client_idが欠けている
		"response_type=code&client_id=00000000-0000-0000-0000-000000000000&redirect_uri=http://example.com&state=xyz",  // scopeが欠けている
		"response_type=code&client_id=00000000-0000-0000-0000-000000000000&scope=read&state=xyz",                       // redirect_uriが欠けている
		"response_type=code&client_id=00000000-0000-0000-0000-000000000000&scope=read&redirect_uri=http://example.com", // stateが欠けている
	}

	for _, params := range invalidParams {
		t.Run("bad request with missing params", func(t *testing.T) {
			t.Parallel()

			r := gin.Default()
			r.LoadHTMLGlob("../../../templates/*")

			handler := &AuthorizeHandler{
				sessionManager: &session.SessionManagerMock{
					NewSessionFunc: func(c *gin.Context) session.SessionClient {
						return &session.SessionClientMock{}
					},
				},
				uc: &AuthorizeUsecaseMock{},
			}

			r.GET("/authorize", handler.Authorize)

			url := "/authorize?" + params

			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
			require.NoError(t, err)

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}

	t.Run("bad request error usecase", func(t *testing.T) {
		t.Parallel()

		r := gin.Default()
		r.LoadHTMLGlob("../../../templates/*")

		handler := &AuthorizeHandler{
			sessionManager: &session.SessionManagerMock{
				NewSessionFunc: func(c *gin.Context) session.SessionClient {
					return &session.SessionClientMock{
						SetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
							return nil
						},
					}
				},
			},
			uc: &AuthorizeUsecaseMock{
				InvokeFunc: func(clientID string, redirectURI string) error {
					return errors.NewUsecaseError(http.StatusBadRequest, "bad requeest")
				},
			},
		}

		r.GET("/authorize", handler.Authorize)

		url := "/authorize?response_type=code&client_id=00000000-0000-0000-0000-000000000000&scope=read&redirect_uri=http://example.com&state=xyz"

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
		require.NoError(t, err)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("internal serer error usecase", func(t *testing.T) {
		t.Parallel()

		r := gin.Default()
		r.LoadHTMLGlob("../../../templates/*")

		handler := &AuthorizeHandler{
			sessionManager: &session.SessionManagerMock{
				NewSessionFunc: func(c *gin.Context) session.SessionClient {
					return &session.SessionClientMock{
						SetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
							return nil
						},
					}
				},
			},
			uc: &AuthorizeUsecaseMock{
				InvokeFunc: func(clientID string, redirectURI string) error {
					return errors.NewUsecaseError(http.StatusInternalServerError, "internal server error")
				},
			},
		}

		r.GET("/authorize", handler.Authorize)

		url := "/authorize?response_type=code&client_id=00000000-0000-0000-0000-000000000000&scope=read&redirect_uri=http://example.com&state=xyz"

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
		require.NoError(t, err)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("internal serer error", func(t *testing.T) {
		t.Parallel()

		r := gin.Default()
		r.LoadHTMLGlob("../../../templates/*")

		handler := &AuthorizeHandler{
			sessionManager: &session.SessionManagerMock{
				NewSessionFunc: func(c *gin.Context) session.SessionClient {
					return &session.SessionClientMock{
						SetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
							return nil
						},
					}
				},
			},
			uc: &AuthorizeUsecaseMock{
				InvokeFunc: func(clientID string, redirectURI string) error {
					return errors.New("internal server error")
				},
			},
		}

		r.GET("/authorize", handler.Authorize)

		url := "/authorize?response_type=code&client_id=00000000-0000-0000-0000-000000000000&scope=read&redirect_uri=http://example.com&state=xyz"

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
		require.NoError(t, err)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("successful authorize", func(t *testing.T) {
		t.Parallel()

		r := gin.Default()
		r.LoadHTMLGlob("../../../templates/*")

		handler := &AuthorizeHandler{
			sessionManager: &session.SessionManagerMock{
				NewSessionFunc: func(c *gin.Context) session.SessionClient {
					return &session.SessionClientMock{
						SetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
							return errors.New("internal server error")
						},
					}
				},
			},
			uc: &AuthorizeUsecaseMock{
				InvokeFunc: func(clientID string, redirectURI string) error {
					return nil
				},
			},
		}

		r.GET("/authorize", handler.Authorize)

		url := "/authorize?response_type=code&client_id=00000000-0000-0000-0000-000000000000&scope=read&redirect_uri=http://example.com&state=xyz"

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
		require.NoError(t, err)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
