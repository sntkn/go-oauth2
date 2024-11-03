package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteToken(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	t.Run("successful delete token", func(t *testing.T) {
		t.Parallel()

		r := gin.Default()

		handler := &DeleteTokenHandler{
			uc: &DeleteTokenUsecaseMock{
				InvokeFunc: func(tokenStr string) error {
					return nil
				},
			},
		}

		r.Use(func(c *gin.Context) {
			c.Set("accessToken", "dummy-token")
		})

		r.DELETE("/token", handler.DeleteToken)

		req, err := http.NewRequestWithContext(context.Background(), http.MethodDelete, "/token", http.NoBody)
		require.NoError(t, err)
		//		req.Header.Add("Authorization", "AccessToken")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("accessToken not exists", func(t *testing.T) {
		t.Parallel()

		r := gin.Default()

		handler := &DeleteTokenHandler{
			uc: &DeleteTokenUsecaseMock{
				InvokeFunc: func(tokenStr string) error {
					return nil
				},
			},
		}
		r.DELETE("/token", handler.DeleteToken)

		req, err := http.NewRequestWithContext(context.Background(), http.MethodDelete, "/token", http.NoBody)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("invalid accessToken type", func(t *testing.T) {
		t.Parallel()

		r := gin.Default()

		handler := &DeleteTokenHandler{
			uc: &DeleteTokenUsecaseMock{
				InvokeFunc: func(tokenStr string) error {
					return errors.NewUsecaseError(http.StatusBadRequest, "bad request")
				},
			},
		}

		r.Use(func(c *gin.Context) {
			c.Set("accessToken", 123)
		})

		r.DELETE("/token", handler.DeleteToken)

		req, err := http.NewRequestWithContext(context.Background(), http.MethodDelete, "/token", http.NoBody)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("internal server error delete token", func(t *testing.T) {
		t.Parallel()

		r := gin.Default()

		handler := &DeleteTokenHandler{
			uc: &DeleteTokenUsecaseMock{
				InvokeFunc: func(tokenStr string) error {
					return errors.New("internal server error")
				},
			},
		}

		r.Use(func(c *gin.Context) {
			c.Set("accessToken", "dummy-token")
		})

		r.DELETE("/token", handler.DeleteToken)

		req, err := http.NewRequestWithContext(context.Background(), http.MethodDelete, "/token", http.NoBody)
		require.NoError(t, err)
		req.Header.Add("Authorization", "AccessToken")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
