package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/common/session"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthorization(t *testing.T) {
	t.Parallel()
	// Ginのテストモードをセット
	gin.SetMode(gin.TestMode)

	cfg, err := config.GetEnv()
	require.NoError(t, err)

	t.Run("successful", func(t *testing.T) {
		t.Parallel()

		// テスト用のルーターを作成
		r := gin.Default()
		r.LoadHTMLGlob("../../../templates/*") // HTMLテンプレートのパスを指定

		handler := &AuthorizationHandler{
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
			uc: &AuthorizationUsecaseMock{
				InvokeFunc: func(input usecases.AuthorizationInput) (string, error) {
					return "https://example.com", nil
				},
			},
			cfg: cfg,
		}

		// authorizationハンドラをセット
		r.POST("/authorization", handler.Authorization)

		// 構造体からフォームデータを生成
		values := url.Values{}
		values.Set("email", "test@example.com")
		values.Add("password", "test1234")

		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/authorization", strings.NewReader(values.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		require.NoError(t, err)

		// レスポンスを記録するためのレスポンスレコーダを作成
		w := httptest.NewRecorder()

		// リクエストをルーターに送信
		r.ServeHTTP(w, req)

		// ステータスコードが302 StatusFoundであることを確認
		assert.Equal(t, http.StatusFound, w.Code)
	})

	t.Run("set session error", func(t *testing.T) {
		t.Parallel()

		// テスト用のルーターを作成
		r := gin.Default()
		r.LoadHTMLGlob("../../../templates/*") // HTMLテンプレートのパスを指定

		handler := &AuthorizationHandler{
			sessionManager: &session.SessionManagerMock{
				NewSessionFunc: func(c *gin.Context) session.SessionClient {
					return &session.SessionClientMock{
						SetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
							return errors.New("set session error")
						},
					}
				},
			},
			uc:  &AuthorizationUsecaseMock{},
			cfg: cfg,
		}

		r.POST("/authorization", handler.Authorization)

		values := url.Values{}
		values.Set("email", "test@example.com")
		values.Add("password", "test1234")

		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/authorization", strings.NewReader(values.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		require.NoError(t, err)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("recirect with invalid param", func(t *testing.T) {
		t.Parallel()
		// テスト用のルーターを作成
		r := gin.Default()
		r.LoadHTMLGlob("../../../templates/*") // HTMLテンプレートのパスを指定

		handler := &AuthorizationHandler{
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
			uc: &AuthorizationUsecaseMock{
				InvokeFunc: func(input usecases.AuthorizationInput) (string, error) {
					return "https://example.com", nil
				},
			},
			cfg: cfg,
		}

		r.POST("/authorization", handler.Authorization)

		values := url.Values{}
		values.Set("email", "test@example.com")

		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/authorization", strings.NewReader(values.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		require.NoError(t, err)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusFound, w.Code)
	})

	t.Run("GetNamedSessionData error", func(t *testing.T) {
		t.Parallel()

		// テスト用のルーターを作成
		r := gin.Default()
		r.LoadHTMLGlob("../../../templates/*") // HTMLテンプレートのパスを指定

		handler := &AuthorizationHandler{
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
							return errors.New("internal server error ")
						},
					}
				},
			},
			uc: &AuthorizationUsecaseMock{
				InvokeFunc: func(input usecases.AuthorizationInput) (string, error) {
					return "https://example.com", nil
				},
			},
			cfg: cfg,
		}

		r.POST("/authorization", handler.Authorization)

		values := url.Values{}
		values.Set("email", "test@example.com")
		values.Add("password", "test1234")

		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/authorization", strings.NewReader(values.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		require.NoError(t, err)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("bad request usecase error with redirect", func(t *testing.T) {
		t.Parallel()

		// テスト用のルーターを作成
		r := gin.Default()
		r.LoadHTMLGlob("../../../templates/*") // HTMLテンプレートのパスを指定

		handler := &AuthorizationHandler{
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
			uc: &AuthorizationUsecaseMock{
				InvokeFunc: func(input usecases.AuthorizationInput) (string, error) {
					return "", errors.NewUsecaseError(http.StatusBadRequest, "bad request")
				},
			},
			cfg: cfg,
		}

		r.POST("/authorization", handler.Authorization)

		values := url.Values{}
		values.Set("email", "test@example.com")
		values.Add("password", "test1234")

		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/authorization", strings.NewReader(values.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		require.NoError(t, err)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusFound, w.Code)
	})

	t.Run("internal server error usecase", func(t *testing.T) {
		t.Parallel()

		// テスト用のルーターを作成
		r := gin.Default()
		r.LoadHTMLGlob("../../../templates/*") // HTMLテンプレートのパスを指定

		handler := &AuthorizationHandler{
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
			uc: &AuthorizationUsecaseMock{
				InvokeFunc: func(input usecases.AuthorizationInput) (string, error) {
					return "", errors.NewUsecaseError(http.StatusInternalServerError, "internal server error")
				},
			},
			cfg: cfg,
		}

		r.POST("/authorization", handler.Authorization)

		values := url.Values{}
		values.Set("email", "test@example.com")
		values.Add("password", "test1234")

		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/authorization", strings.NewReader(values.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		require.NoError(t, err)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("internal server error", func(t *testing.T) {
		t.Parallel()

		// テスト用のルーターを作成
		r := gin.Default()
		r.LoadHTMLGlob("../../../templates/*") // HTMLテンプレートのパスを指定

		handler := &AuthorizationHandler{
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
			uc: &AuthorizationUsecaseMock{
				InvokeFunc: func(input usecases.AuthorizationInput) (string, error) {
					return "", errors.New("internal server error")
				},
			},
			cfg: cfg,
		}

		r.POST("/authorization", handler.Authorization)

		values := url.Values{}
		values.Set("email", "test@example.com")
		values.Add("password", "test1234")

		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/authorization", strings.NewReader(values.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		require.NoError(t, err)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("del session error", func(t *testing.T) {
		t.Parallel()

		r := gin.Default()
		r.LoadHTMLGlob("../../../templates/*")

		handler := &AuthorizationHandler{
			sessionManager: &session.SessionManagerMock{
				NewSessionFunc: func(c *gin.Context) session.SessionClient {
					return &session.SessionClientMock{
						SetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
							return nil
						},
						DelSessionDataFunc: func(c *gin.Context, key string) error {
							return errors.New("internal server error")
						},
						GetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
							return nil
						},
					}
				},
			},
			uc: &AuthorizationUsecaseMock{
				InvokeFunc: func(input usecases.AuthorizationInput) (string, error) {
					return "https://example.com", nil
				},
			},
			cfg: cfg,
		}

		r.POST("/authorization", handler.Authorization)

		values := url.Values{}
		values.Set("email", "test@example.com")
		values.Add("password", "test1234")

		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/authorization", strings.NewReader(values.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		require.NoError(t, err)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
