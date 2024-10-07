package user

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/internal/accesstoken"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUser(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	t.Run("successful get user", func(t *testing.T) {
		t.Parallel()

		r := gin.Default()

		token, err := accesstoken.Generate(accesstoken.TokenParams{
			UserID:    uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			ClientID:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			Scope:     "read",
			ExpiresAt: time.Now().Add(1 * time.Hour),
		})
		require.NoError(t, err)

		handler := &GetUserHandler{
			uc: &GetUserUsecaseMock{
				InvokeFunc: func(userID uuid.UUID) (repository.User, error) {
					return repository.User{}, nil
				},
			},
		}
		r.GET("/me", handler.GetUser)

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/me", http.NoBody)
		require.NoError(t, err)
		req.Header.Add("Authorization", token)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad request get user", func(t *testing.T) {
		t.Parallel()

		r := gin.Default()

		handler := &GetUserHandler{
			uc: &GetUserUsecaseMock{
				InvokeFunc: func(userID uuid.UUID) (repository.User, error) {
					return repository.User{}, nil
				},
			},
		}
		r.GET("/me", handler.GetUser)

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/me", http.NoBody)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("usecase error get user", func(t *testing.T) {
		t.Parallel()

		r := gin.Default()

		token, err := accesstoken.Generate(accesstoken.TokenParams{
			UserID:    uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			ClientID:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			Scope:     "read",
			ExpiresAt: time.Now().Add(1 * time.Hour),
		})
		require.NoError(t, err)

		handler := &GetUserHandler{
			uc: &GetUserUsecaseMock{
				InvokeFunc: func(userID uuid.UUID) (repository.User, error) {
					return repository.User{}, errors.NewUsecaseError(http.StatusBadRequest, "bad request")
				},
			},
		}
		r.GET("/me", handler.GetUser)

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/me", http.NoBody)
		require.NoError(t, err)
		req.Header.Add("Authorization", token)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("internal server error get user", func(t *testing.T) {
		t.Parallel()

		r := gin.Default()

		token, err := accesstoken.Generate(accesstoken.TokenParams{
			UserID:    uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			ClientID:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			Scope:     "read",
			ExpiresAt: time.Now().Add(1 * time.Hour),
		})
		require.NoError(t, err)

		handler := &GetUserHandler{
			uc: &GetUserUsecaseMock{
				InvokeFunc: func(userID uuid.UUID) (repository.User, error) {
					return repository.User{}, errors.New("internal server error")
				},
			},
		}
		r.GET("/me", handler.GetUser)

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/me", http.NoBody)
		require.NoError(t, err)
		req.Header.Add("Authorization", token)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
