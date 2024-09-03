package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/internal/accesstoken"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type MockGetUserUsecase struct {
	mock.Mock
}

func (m *MockGetUserUsecase) Invoke(c *gin.Context) (repository.User, error) {
	args := m.Called(c)
	return args.Get(0).(repository.User), args.Error(1)
}

func TestGetUserHandler(t *testing.T) {
	t.Parallel()
	// Ginのテストモードをセット
	gin.SetMode(gin.TestMode)

	// テスト用のルーターを作成
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("db", &repository.OAuth2RepositoryMock{
			FindUserFunc: func(id uuid.UUID) (repository.User, error) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("test1234"), bcrypt.DefaultCost)
				return repository.User{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					Name:      "client Name",
					Email:     "test@example.com",
					Password:  string(hashedPassword),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
		})
		c.Next() // 次のミドルウェア/ハンドラへ
	})

	// サインインハンドラをセット
	r.GET("/me", func(c *gin.Context) {
		GetUserHandler(c)
	})

	token, err := accesstoken.Generate(accesstoken.TokenParams{
		UserID:    uuid.MustParse("00000000-0000-0000-0000-000000000000"),
		ClientID:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
		Scope:     "read",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	})
	require.NoError(t, err)

	// テスト用のHTTPリクエストとレスポンスレコーダを作成
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/me", http.NoBody)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	require.NoError(t, err)

	// レスポンスを記録するためのレスポンスレコーダを作成
	w := httptest.NewRecorder()

	// リクエストをルーターに送信
	r.ServeHTTP(w, req)

	// ステータスコードが200 OKであることを確認
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUser(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	t.Run("successful get user", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		mockUC := new(MockGetUserUsecase)
		mockUC.On("Invoke", mock.Anything).Return(repository.User{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			Name:      "client Name",
			Email:     "test@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil)

		getUser(c, mockUC)

		assert.Equal(t, http.StatusOK, w.Code)

		mockUC.AssertExpectations(t)
	})

	t.Run("unauthorized error", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		mockUC := new(MockGetUserUsecase)
		mockUC.On("Invoke", mock.Anything).Return(repository.User{}, &cerrs.UsecaseError{Code: http.StatusUnauthorized})

		getUser(c, mockUC)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		mockUC := new(MockGetUserUsecase)
		mockUC.On("Invoke", mock.Anything).Return(repository.User{}, errors.New("internal error"))

		getUser(c, mockUC)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockUC.AssertExpectations(t)
	})
}
