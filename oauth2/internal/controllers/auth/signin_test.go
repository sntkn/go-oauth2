package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/flashmessage"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSigninSuccessful(t *testing.T) {
	t.Parallel()
	// Ginのテストモードをセット
	gin.SetMode(gin.TestMode)

	// テスト用のルーターを作成
	r := gin.Default()
	r.LoadHTMLGlob("../../../templates/*") // HTMLテンプレートのパスを指定

	handler := &SigninHandler{
		sessionManager: &session.SessionManagerMock{
			NewSessionFunc: func(c *gin.Context) session.SessionClient {
				return &session.SessionClientMock{
					GetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
						switch v := t.(type) {
						case *AuthorizeInput:
							*v = AuthorizeInput{ClientID: "1234-abcd-qwer-asdf"}
						case *flashmessage.Messages:
							*v = flashmessage.Messages{}
						default:
							return errors.New("interface conversion error")
						}
						return nil
					},
					SetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
						return nil
					},
					FlushNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
						return nil
					},
				}
			},
		},
	}

	// サインインハンドラをセット
	r.GET("/signin", handler.Signin)

	// テスト用のHTTPリクエストとレスポンスレコーダを作成
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/signin", http.NoBody)
	require.NoError(t, err)

	// レスポンスを記録するためのレスポンスレコーダを作成
	w := httptest.NewRecorder()

	// リクエストをルーターに送信
	r.ServeHTTP(w, req)

	// ステータスコードが200 OKであることを確認
	assert.Equal(t, http.StatusOK, w.Code)

	// レスポンスに含まれるべきHTMLコンテンツが含まれているか確認
	assert.Contains(t, w.Body.String(), "<h2>Signin</h2>")
}

func TestSigninClientNotFound(t *testing.T) {
	t.Parallel()
	// Ginのテストモードをセット
	gin.SetMode(gin.TestMode)

	// テスト用のルーターを作成
	r := gin.Default()
	r.LoadHTMLGlob("../../../templates/*") // HTMLテンプレートのパスを指定

	handler := &SigninHandler{
		sessionManager: &session.SessionManagerMock{
			NewSessionFunc: func(c *gin.Context) session.SessionClient {
				return &session.SessionClientMock{
					GetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
						switch v := t.(type) {
						case *AuthorizeInput:
							*v = AuthorizeInput{}
						case *flashmessage.Messages:
							*v = flashmessage.Messages{}
						default:
							return errors.New("interface conversion error")
						}
						return nil
					},
					SetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
						return nil
					},
					FlushNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
						return nil
					},
				}
			},
		},
	}

	r.GET("/signin", handler.Signin)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/signin", http.NoBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
