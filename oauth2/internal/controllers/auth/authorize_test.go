package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthorizeUsecase struct {
	mock.Mock
}

func (m *MockAuthorizeUsecase) Invoke(c *gin.Context, ClientID string, redirectURI string) error {
	args := m.Called(c)
	return args.Error(0)
}

func TestAuthorizeHandler(t *testing.T) {
	// Ginのテストモードをセット
	gin.SetMode(gin.TestMode)

	// テスト用のルーターを作成
	r := gin.Default()
	r.LoadHTMLGlob("../../../templates/*") // HTMLテンプレートのパスを指定

	r.Use(func(c *gin.Context) {
		c.Set("db", &repository.OAuth2RepositoryMock{
			FindClientByClientIDFunc: func(clientID string) (repository.Client, error) {
				return repository.Client{
					ID:           uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					Name:         "client Name",
					RedirectURIs: "http://example.com",
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				}, nil
			},
		})
		c.Set("session", &session.SessionClientMock{
			SetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				return nil
			},
		})
		c.Set("cfg", &config.Config{})
		c.Next() // 次のミドルウェア/ハンドラへ
	})

	// authorizeハンドラをセット
	r.GET("/authorize", func(c *gin.Context) {
		AuthorizeHandler(c)
	})

	// テスト用のHTTPリクエストとレスポンスレコーダを作成
	url := "/authorize?response_type=code&client_id=00000000-0000-0000-0000-000000000000&scope=read&redirect_uri=http://example.com&state=xyz"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	assert.NoError(t, err)

	// レスポンスを記録するためのレスポンスレコーダを作成
	w := httptest.NewRecorder()

	// リクエストをルーターに送信
	r.ServeHTTP(w, req)

	// ステータスコードが302 StatusFoundであることを確認
	assert.Equal(t, http.StatusFound, w.Code)
}

func TestAuthorize(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful authorize", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.LoadHTMLGlob("../../../templates/*")
		c.Request = httptest.NewRequest(http.MethodGet, "/?response_type=code&client_id=00000000-0000-0000-0000-000000000000&scope=read&redirect_uri=http://example.com&state=xyz", nil)

		s := &session.SessionClientMock{
			SetNamedSessionDataFunc: func(c *gin.Context, key string, t any) error {
				return nil
			},
		}
		mockUC := new(MockAuthorizeUsecase)
		mockUC.On("Invoke", mock.Anything).Return(nil)

		authorize(c, s, mockUC)

		assert.Equal(t, http.StatusFound, w.Code)

		mockUC.AssertExpectations(t)
	})

	t.Run("bad request error(validation)", func(t *testing.T) {

		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.LoadHTMLGlob("../../../templates/*")
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

		s := &session.SessionClientMock{}
		mockUC := new(MockAuthorizeUsecase)

		authorize(c, s, mockUC)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("bad request error(usecase)", func(t *testing.T) {

		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.LoadHTMLGlob("../../../templates/*")
		c.Request = httptest.NewRequest(http.MethodGet, "/?response_type=code&client_id=00000000-0000-0000-0000-000000000000&scope=read&redirect_uri=http://example.com&state=xyz", nil)

		s := &session.SessionClientMock{}
		mockUC := new(MockAuthorizeUsecase)
		mockUC.On("Invoke", mock.Anything).Return(&cerrs.UsecaseError{Code: http.StatusBadRequest})

		authorize(c, s, mockUC)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {

		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		r.LoadHTMLGlob("../../../templates/*")
		c.Request = httptest.NewRequest(http.MethodGet, "/?response_type=code&client_id=00000000-0000-0000-0000-000000000000&scope=read&redirect_uri=http://example.com&state=xyz", nil)

		s := &session.SessionClientMock{}
		mockUC := new(MockAuthorizeUsecase)
		mockUC.On("Invoke", mock.Anything).Return(errors.New("internal error"))

		authorize(c, s, mockUC)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockUC.AssertExpectations(t)
	})
}
