package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/entity"
	"github.com/sntkn/go-oauth2/oauth2/internal/flashmessage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	// Ginのテストモードをセット
	gin.SetMode(gin.TestMode)

	t.Run("signup finished successful", func(t *testing.T) {
		t.Parallel()
		// テスト用のルーターを作成
		r := gin.Default()
		r.LoadHTMLGlob("../../../templates/*") // HTMLテンプレートのパスを指定

		r.Use(func(c *gin.Context) {
			c.Set("flashMessages", &flashmessage.Messages{})
			c.Next() // 次のミドルウェア/ハンドラへ
		})

		// サインインハンドラをセット
		r.GET("/signup_finished", func(c *gin.Context) {
			SignupFinishedHandler(c)
		})

		// テスト用のHTTPリクエストとレスポンスレコーダを作成
		req, err := http.NewRequest(http.MethodGet, "/signup_finished", nil)
		assert.NoError(t, err)

		// レスポンスを記録するためのレスポンスレコーダを作成
		w := httptest.NewRecorder()

		// リクエストをルーターに送信
		r.ServeHTTP(w, req)

		// ステータスコードが200 OKであることを確認
		assert.Equal(t, http.StatusOK, w.Code)

		// レスポンスに含まれるべきHTMLコンテンツが含まれているか確認
		assert.Contains(t, w.Body.String(), "User creation was successful.")
	})
}
