package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/internal/entity"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/validation"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCreateTokenByCodeUsecase struct {
	mock.Mock
}

func (m *MockCreateTokenByCodeUsecase) Invoke(c *gin.Context, authCode string) (entity.AuthTokens, error) {
	args := m.Called(c)
	return args.Get(0).(entity.AuthTokens), args.Error(1)
}

type MockCreateTokenByRefreshTokenUsecase struct {
	mock.Mock
}

func (m *MockCreateTokenByRefreshTokenUsecase) Invoke(c *gin.Context, refreshToken string) (entity.AuthTokens, error) {
	args := m.Called(c)
	return args.Get(0).(entity.AuthTokens), args.Error(1)
}

func TestCreateTokenHandler(t *testing.T) {
	t.Parallel()
	// Ginのテストモードをセット
	gin.SetMode(gin.TestMode)

	t.Run("successful create token by refresh token", func(t *testing.T) {
		t.Parallel()

		// テスト用のルーターを作成
		r := gin.Default()

		r.Use(func(c *gin.Context) {
			c.Set("db", &repository.OAuth2RepositoryMock{
				FindValidOAuth2CodeFunc: func(code string, expiresAt time.Time) (repository.Code, error) {
					return repository.Code{
						UserID:    uuid.MustParse("00000000-0000-0000-0000-000000000000"),
						ClientID:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
						Scope:     "",
						ExpiresAt: time.Now().Add(1 * time.Hour),
					}, nil
				},
				RegisterTokenFunc: func(t *repository.Token) error {
					return nil
				},
				RegisterRefreshTokenFunc: func(t *repository.RefreshToken) error {
					return nil
				},
				RevokeCodeFunc: func(code string) error {
					return nil
				},
			})
			c.Set("cfg", &config.Config{})
			c.Next() // 次のミドルウェア/ハンドラへ
		})

		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			v.RegisterValidation("required_with_field_value", validation.RequiredWithFieldValue)
		}

		// サインインハンドラをセット
		r.POST("/token", func(c *gin.Context) {
			CreateTokenHandler(c)
		})

		// テスト用のHTTPリクエストとレスポンスレコーダを作成
		// Create an example user for testing
		exampleToken := TokenInput{
			GrantType: "authorization_code",
			Code:      "test",
		}
		j, _ := json.Marshal(exampleToken)
		req, err := http.NewRequest(http.MethodPost, "/token", strings.NewReader(string(j)))
		req.Header.Add("Authorization", "AccessToken")
		req.Header.Set("Content-Type", "application/json")

		assert.NoError(t, err)

		// レスポンスを記録するためのレスポンスレコーダを作成
		w := httptest.NewRecorder()

		// リクエストをルーターに送信
		r.ServeHTTP(w, req)

		// ステータスコードが200 OKであることを確認
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("successful create token by code", func(t *testing.T) {
		t.Parallel()
		// テスト用のルーターを作成
		r := gin.Default()

		r.Use(func(c *gin.Context) {
			c.Set("db", &repository.OAuth2RepositoryMock{
				FindValidRefreshTokenFunc: func(refreshToken string, expiresAt time.Time) (repository.RefreshToken, error) {
					return repository.RefreshToken{
						RefreshToken: "test",
						AccessToken:  "test",
						ExpiresAt:    time.Now().Add(1 * time.Hour),
						CreatedAt:    time.Now(),
						UpdatedAt:    time.Now(),
					}, nil
				},
				FindTokenFunc: func(accessToken string) (repository.Token, error) {
					return repository.Token{
						AccessToken: "test",
						UserID:      uuid.MustParse("00000000-0000-0000-0000-000000000000"),
						ClientID:    uuid.MustParse("00000000-0000-0000-0000-000000000000"),
						Scope:       "",
						ExpiresAt:   time.Now().Add(1 * time.Hour),
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					}, nil
				},
				RegisterTokenFunc: func(t *repository.Token) error {
					return nil
				},
				RegisterRefreshTokenFunc: func(t *repository.RefreshToken) error {
					return nil
				},
				RevokeTokenFunc: func(accessToken string) error {
					return nil
				},
				RevokeRefreshTokenFunc: func(refreshToken string) error {
					return nil
				},
			})
			c.Set("cfg", &config.Config{})
			c.Next() // 次のミドルウェア/ハンドラへ
		})

		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			v.RegisterValidation("required_with_field_value", validation.RequiredWithFieldValue)
		}

		// サインインハンドラをセット
		r.POST("/token", func(c *gin.Context) {
			CreateTokenHandler(c)
		})

		// テスト用のHTTPリクエストとレスポンスレコーダを作成
		// Create an example user for testing
		exampleToken := TokenInput{
			GrantType:    "refresh_token",
			RefreshToken: "test",
		}
		j, _ := json.Marshal(exampleToken)
		req, err := http.NewRequest(http.MethodPost, "/token", strings.NewReader(string(j)))
		req.Header.Add("Authorization", "AccessToken")
		req.Header.Set("Content-Type", "application/json")

		assert.NoError(t, err)

		// レスポンスを記録するためのレスポンスレコーダを作成
		w := httptest.NewRecorder()

		// リクエストをルーターに送信
		r.ServeHTTP(w, req)

		// ステータスコードが200 OKであることを確認
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestCreateToken(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	t.Run("successful sign-in", func(t *testing.T) {
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
