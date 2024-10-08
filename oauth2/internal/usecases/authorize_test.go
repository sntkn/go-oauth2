package usecases

import (
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthorizeInvoke(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	t.Run("successful invoke", func(t *testing.T) {
		t.Parallel()
		clientID := "00000000-0000-0000-0000-000000000000"
		redirectURI := "http://example.com"

		db := &repository.OAuth2RepositoryMock{
			FindClientByClientIDFunc: func(clientID string) (repository.Client, error) {
				return repository.Client{
					ID:           uuid.MustParse(clientID),
					Name:         "client Name",
					RedirectURIs: redirectURI,
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				}, nil
			},
		}

		authorize := NewAuthorize(db)
		err := authorize.Invoke(clientID, redirectURI)
		require.NoError(t, err)
	})

	t.Run("missing to find client_id", func(t *testing.T) {
		t.Parallel()
		clientID := "00000000-0000-0000-0000-000000000000"
		redirectURI := "http://example.com"

		db := &repository.OAuth2RepositoryMock{
			FindClientByClientIDFunc: func(clientID string) (repository.Client, error) {
				return repository.Client{}, &errors.UsecaseError{Code: http.StatusBadRequest}
			},
		}

		authorize := NewAuthorize(db)
		err := authorize.Invoke(clientID, redirectURI)
		require.Error(t, err)
		assert.IsType(t, &errors.UsecaseError{}, err)
	})

	t.Run("database error", func(t *testing.T) {
		t.Parallel()
		clientID := "00000000-0000-0000-0000-000000000000"
		redirectURI := "http://example.com"

		db := &repository.OAuth2RepositoryMock{
			FindClientByClientIDFunc: func(clientID string) (repository.Client, error) {
				return repository.Client{}, errors.New("internal error")
			},
		}

		authorize := NewAuthorize(db)
		err := authorize.Invoke(clientID, redirectURI)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "internal error")
	})

	t.Run("redirect_uri does not match", func(t *testing.T) {
		t.Parallel()

		db := &repository.OAuth2RepositoryMock{
			FindClientByClientIDFunc: func(clientID string) (repository.Client, error) {
				return repository.Client{
					ID:           uuid.MustParse(clientID),
					Name:         "client Name",
					RedirectURIs: "http://example.com",
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				}, nil
			},
		}
		clientID := "00000000-0000-0000-0000-000000000000"
		redirectURI := "http://example.com/not/match"

		authorize := NewAuthorize(db)
		err := authorize.Invoke(clientID, redirectURI)
		require.Error(t, err)
		assert.IsType(t, &errors.UsecaseError{}, err)
		assert.Contains(t, err.Error(), "redirect uri does not match")
	})
}
