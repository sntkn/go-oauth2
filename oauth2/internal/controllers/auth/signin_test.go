package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/flashmessage"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestSigninHandler(t *testing.T) {
	// Ginのテストモードをセット
	gin.SetMode(gin.TestMode)

	// テスト用のルーターを作成
	r := gin.Default()
	r.LoadHTMLGlob("../../../templates/*") // HTMLテンプレートのパスを指定

	r.Use(func(c *gin.Context) {
		c.Set("session", &session.SessionClientMock{
			GetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				*t.(*usecases.AuthorizeInput) = usecases.AuthorizeInput{ClientID: "1234-abcd-qwer-asdf"}
				return nil
			},
			FlushNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				return nil
			},
		})
		c.Set("flashMessages", &flashmessage.Messages{})
		c.Set("cfg", &config.Config{})
		c.Next() // 次のミドルウェア/ハンドラへ
	})

	// サインインハンドラをセット
	r.GET("/signin", func(c *gin.Context) {
		SigninHandler(c)
	})

	// テスト用のHTTPリクエストとレスポンスレコーダを作成
	req, err := http.NewRequest(http.MethodGet, "/signin", nil)
	assert.NoError(t, err)

	// レスポンスを記録するためのレスポンスレコーダを作成
	w := httptest.NewRecorder()

	// リクエストをルーターに送信
	r.ServeHTTP(w, req)

	// ステータスコードが200 OKであることを確認
	assert.Equal(t, http.StatusOK, w.Code)

	// レスポンスに含まれるべきHTMLコンテンツが含まれているか確認
	assert.Contains(t, w.Body.String(), "<h2>Signin</h2>")
}
