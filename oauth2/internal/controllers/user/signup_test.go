package user

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/entity"
	"github.com/sntkn/go-oauth2/oauth2/internal/flashmessage"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
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
	// Ginのテストモードをセット
	gin.SetMode(gin.TestMode)

	// テスト用のルーターを作成
	r := gin.Default()
	r.LoadHTMLGlob("../../../templates/*") // HTMLテンプレートのパスを指定
	form := entity.SessionRegistrationForm{
		Name:  "test",
		Email: "test@example.com",
		Error: "",
	}

	r.Use(func(c *gin.Context) {
		c.Set("session", &session.SessionClientMock{
			FlushNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				*t.(*entity.SessionRegistrationForm) = form
				return nil
			},
		})
		c.Set("flashMessages", &flashmessage.Messages{})
		c.Set("cfg", &config.Config{})
		c.Next() // 次のミドルウェア/ハンドラへ
	})

	// サインインハンドラをセット
	r.GET("/signup", func(c *gin.Context) {
		SignupHandler(c)
	})

	// テスト用のHTTPリクエストとレスポンスレコーダを作成
	req, err := http.NewRequest(http.MethodGet, "/signup", http.NoBody)
	require.NoError(t, err)

	// レスポンスを記録するためのレスポンスレコーダを作成
	w := httptest.NewRecorder()

	// リクエストをルーターに送信
	r.ServeHTTP(w, req)

	// ステータスコードが200 OKであることを確認
	assert.Equal(t, http.StatusOK, w.Code)

	// レスポンスに含まれるべきHTMLコンテンツが含まれているか確認
	assert.Contains(t, w.Body.String(), "<h2>Signup</h2>")
}

func TestSignup(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	t.Run("successful signup", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.LoadHTMLGlob("../../../templates/*")
		form := entity.SessionRegistrationForm{
			Name:  "test",
			Email: "test@example.com",
			Error: "",
		}

		mess := &flashmessage.Messages{}
		mockUC := new(MockSignupUsecase)
		mockUC.On("Invoke", mock.Anything).Return(form, nil)

		signup(c, mess, mockUC)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "<h2>Signup</h2>")

		mockUC.AssertExpectations(t)
	})

	t.Run("usecase internal server error", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.LoadHTMLGlob("../../../templates/*")

		mess := &flashmessage.Messages{}
		mockUC := new(MockSignupUsecase)
		mockUC.On("Invoke", mock.Anything).Return(entity.SessionRegistrationForm{}, &cerrs.UsecaseError{Code: http.StatusInternalServerError})

		signup(c, mess, mockUC)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.LoadHTMLGlob("../../../templates/*")

		mess := &flashmessage.Messages{}
		mockUC := new(MockSignupUsecase)
		mockUC.On("Invoke", mock.Anything).Return(entity.SessionRegistrationForm{}, errors.New("internal error"))

		signup(c, mess, mockUC)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockUC.AssertExpectations(t)
	})
}
