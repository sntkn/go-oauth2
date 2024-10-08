package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/sntkn/go-oauth2/oauth2/internal/entity"
	"github.com/sntkn/go-oauth2/oauth2/internal/validation"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTokenHandler(t *testing.T) {
	t.Parallel()
	// Ginのテストモードをセット
	gin.SetMode(gin.TestMode)

	// ginのバリデーションエンジンにカスタムバリデーションを登録する
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("required_with_field_value", validation.RequiredWithFieldValue)
		require.NoError(t, err)
	}

	t.Run("successful create token by auth code", func(t *testing.T) {
		t.Parallel()

		r := gin.Default()

		handler := &CreateTokenHandler{
			tokenUC: &CreateTokenByCodeUsecaseMock{
				InvokeFunc: func(refreshToken string) (*entity.AuthTokens, error) {
					return &entity.AuthTokens{
						AccessToken:  "test",
						RefreshToken: "test",
						Expiry:       time.Now().Add(1 * time.Hour).Unix(),
					}, nil
				},
			},
			refreshTokenUC: &CreateTokenByRefreshTokenUsecaseMock{},
		}

		r.POST("/token", handler.CreateToken)

		exampleToken := TokenInput{
			GrantType: "authorization_code",
			Code:      "test",
		}
		j, err := json.Marshal(exampleToken)
		require.NoError(t, err)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/token", strings.NewReader(string(j)))
		require.NoError(t, err)
		req.Header.Add("Authorization", "AccessToken")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("successful create token by refresh token", func(t *testing.T) {
		t.Parallel()

		r := gin.Default()

		handler := &CreateTokenHandler{
			tokenUC: &CreateTokenByCodeUsecaseMock{},
			refreshTokenUC: &CreateTokenByRefreshTokenUsecaseMock{
				InvokeFunc: func(refreshToken string) (*entity.AuthTokens, error) {
					return &entity.AuthTokens{
						AccessToken:  "test",
						RefreshToken: "test",
						Expiry:       time.Now().Add(1 * time.Hour).Unix(),
					}, nil
				},
			},
		}

		r.POST("/token", handler.CreateToken)

		exampleToken := TokenInput{
			GrantType:    "refresh_token",
			RefreshToken: "test",
		}
		j, err := json.Marshal(exampleToken)
		require.NoError(t, err)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/token", strings.NewReader(string(j)))
		require.NoError(t, err)
		req.Header.Add("Authorization", "AccessToken")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid grant type", func(t *testing.T) {
		t.Parallel()

		r := gin.Default()

		handler := &CreateTokenHandler{
			tokenUC:        &CreateTokenByCodeUsecaseMock{},
			refreshTokenUC: &CreateTokenByRefreshTokenUsecaseMock{},
		}

		r.POST("/token", handler.CreateToken)

		exampleToken := TokenInput{
			GrantType:    "invalid_type",
			RefreshToken: "test",
			Code:         "test",
		}
		j, err := json.Marshal(exampleToken)
		require.NoError(t, err)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/token", strings.NewReader(string(j)))
		require.NoError(t, err)
		req.Header.Add("Authorization", "AccessToken")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase error create token by auth code", func(t *testing.T) {
		t.Parallel()

		r := gin.Default()

		handler := &CreateTokenHandler{
			tokenUC: &CreateTokenByCodeUsecaseMock{
				InvokeFunc: func(refreshToken string) (*entity.AuthTokens, error) {
					return nil, errors.NewUsecaseError(http.StatusBadRequest, "bad request")
				},
			},
			refreshTokenUC: &CreateTokenByRefreshTokenUsecaseMock{},
		}

		r.POST("/token", handler.CreateToken)

		exampleToken := TokenInput{
			GrantType: "authorization_code",
			Code:      "test",
		}
		j, err := json.Marshal(exampleToken)
		require.NoError(t, err)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/token", strings.NewReader(string(j)))
		require.NoError(t, err)
		req.Header.Add("Authorization", "AccessToken")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("internal server error create token by auth code", func(t *testing.T) {
		t.Parallel()

		r := gin.Default()

		handler := &CreateTokenHandler{
			tokenUC: &CreateTokenByCodeUsecaseMock{
				InvokeFunc: func(refreshToken string) (*entity.AuthTokens, error) {
					return nil, errors.New("internal server error")
				},
			},
			refreshTokenUC: &CreateTokenByRefreshTokenUsecaseMock{},
		}

		r.POST("/token", handler.CreateToken)

		exampleToken := TokenInput{
			GrantType: "authorization_code",
			Code:      "test",
		}
		j, err := json.Marshal(exampleToken)
		require.NoError(t, err)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/token", strings.NewReader(string(j)))
		require.NoError(t, err)
		req.Header.Add("Authorization", "AccessToken")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("usecase error create token by refresh token", func(t *testing.T) {
		t.Parallel()

		r := gin.Default()

		handler := &CreateTokenHandler{
			tokenUC: &CreateTokenByCodeUsecaseMock{},
			refreshTokenUC: &CreateTokenByRefreshTokenUsecaseMock{
				InvokeFunc: func(refreshToken string) (*entity.AuthTokens, error) {
					return nil, errors.NewUsecaseError(http.StatusBadRequest, "bad request")
				},
			},
		}

		r.POST("/token", handler.CreateToken)

		exampleToken := TokenInput{
			GrantType:    "refresh_token",
			RefreshToken: "test",
		}
		j, err := json.Marshal(exampleToken)
		require.NoError(t, err)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/token", strings.NewReader(string(j)))
		require.NoError(t, err)
		req.Header.Add("Authorization", "AccessToken")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("internal server error create token by refresh token", func(t *testing.T) {
		t.Parallel()

		r := gin.Default()

		handler := &CreateTokenHandler{
			tokenUC: &CreateTokenByCodeUsecaseMock{},
			refreshTokenUC: &CreateTokenByRefreshTokenUsecaseMock{
				InvokeFunc: func(refreshToken string) (*entity.AuthTokens, error) {
					return nil, errors.New("internal server error")
				},
			},
		}

		r.POST("/token", handler.CreateToken)

		exampleToken := TokenInput{
			GrantType:    "refresh_token",
			RefreshToken: "test",
		}
		j, err := json.Marshal(exampleToken)
		require.NoError(t, err)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/token", strings.NewReader(string(j)))
		require.NoError(t, err)
		req.Header.Add("Authorization", "AccessToken")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
