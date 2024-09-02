package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/entity"
	"github.com/sntkn/go-oauth2/oauth2/internal/flashmessage"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSigninUsecase struct {
	mock.Mock
}

func (m *MockSigninUsecase) Invoke(c *gin.Context) (entity.SessionSigninForm, error) {
	args := m.Called(c)
	return args.Get(0).(entity.SessionSigninForm), args.Error(1)
}

func TestSigninHandler(t *testing.T) {
	t.Parallel()
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
	req, err := http.NewRequest(http.MethodGet, "/signin", http.NoBody)
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

func TestSignin(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	t.Run("successful sign-in", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.LoadHTMLGlob("../../../templates/*")

		mess := &flashmessage.Messages{}
		mockUC := new(MockSigninUsecase)
		mockForm := entity.SessionSigninForm{}
		mockUC.On("Invoke", mock.Anything).Return(mockForm, nil)

		signin(c, mess, mockUC)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "<h2>Signin</h2>")

		mockUC.AssertExpectations(t)
	})

	t.Run("bad request error", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.LoadHTMLGlob("../../../templates/*")

		mess := &flashmessage.Messages{}
		mockUC := new(MockSigninUsecase)
		mockUC.On("Invoke", mock.Anything).Return(entity.SessionSigninForm{}, &cerrs.UsecaseError{Code: http.StatusBadRequest})

		signin(c, mess, mockUC)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.LoadHTMLGlob("../../../templates/*")

		mess := &flashmessage.Messages{}
		mockUC := new(MockSigninUsecase)
		mockUC.On("Invoke", mock.Anything).Return(entity.SessionSigninForm{}, errors.New("internal error"))

		signin(c, mess, mockUC)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockUC.AssertExpectations(t)
	})
}
