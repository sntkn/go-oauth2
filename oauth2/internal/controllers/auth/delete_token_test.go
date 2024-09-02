package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockDeleteTokenUsecase struct {
	mock.Mock
}

func (m *MockDeleteTokenUsecase) Invoke(c *gin.Context) error {
	args := m.Called(c)
	return args.Error(0)
}

func TestDeleteTokenHandler(t *testing.T) {
	t.Parallel()
	// Ginのテストモードをセット
	gin.SetMode(gin.TestMode)

	// テスト用のルーターを作成
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("db", &repository.OAuth2RepositoryMock{
			RevokeTokenFunc: func(accessToken string) error {
				return nil
			},
		})
		c.Next() // 次のミドルウェア/ハンドラへ
	})

	// サインインハンドラをセット
	r.GET("/delete_token", func(c *gin.Context) {
		DeleteTokenHandler(c)
	})

	// テスト用のHTTPリクエストとレスポンスレコーダを作成
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/delete_token", http.NoBody)
	require.NoError(t, err)
	req.Header.Add("Authorization", "AccessToken")
	require.NoError(t, err)

	// レスポンスを記録するためのレスポンスレコーダを作成
	w := httptest.NewRecorder()

	// リクエストをルーターに送信
	r.ServeHTTP(w, req)

	// ステータスコードが200 OKであることを確認
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteToken(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	t.Run("successful sign-in", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		mockUC := new(MockDeleteTokenUsecase)
		mockUC.On("Invoke", mock.Anything).Return(nil)

		deleteToken(c, mockUC)

		assert.Equal(t, http.StatusOK, w.Code)

		mockUC.AssertExpectations(t)
	})

	t.Run("bad request error", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		mockUC := new(MockDeleteTokenUsecase)
		mockUC.On("Invoke", mock.Anything).Return(&cerrs.UsecaseError{Code: http.StatusBadRequest})

		deleteToken(c, mockUC)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		t.Parallel()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		mockUC := new(MockDeleteTokenUsecase)
		mockUC.On("Invoke", mock.Anything).Return(errors.New("internal error"))

		deleteToken(c, mockUC)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockUC.AssertExpectations(t)
	})
}
