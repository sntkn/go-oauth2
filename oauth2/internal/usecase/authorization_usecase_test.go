package usecase

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/domain"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConsent_Success(t *testing.T) {
	mockClientRepo := &domain.ClientRepositoryMock{
		FindClientByClientIDFunc: func(clientID uuid.UUID) (domain.Client, error) {
			return &domain.ClientMock{
				IsNotFoundFunc: func() bool {
					return false
				},
			}, nil
		},
	}

	uc := NewAuthorizationUsecase(mockClientRepo, nil, nil, nil, nil)
	_, err := uc.Consent(uuid.New())
	require.NoError(t, err)
}

func TestConsent_FindClientError(t *testing.T) {
	mockClientRepo := &domain.ClientRepositoryMock{
		FindClientByClientIDFunc: func(clientID uuid.UUID) (domain.Client, error) {
			return nil, errors.New("FindClientByClientID error")
		},
	}

	uc := NewAuthorizationUsecase(mockClientRepo, nil, nil, nil, nil)
	_, err := uc.Consent(uuid.New())
	require.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*errors.UsecaseError).Code)
	assert.Equal(t, "FindClientByClientID error", err.(*errors.UsecaseError).Message)
}

func TestConsent_ClientNotFound(t *testing.T) {
	mockClientRepo := &domain.ClientRepositoryMock{
		FindClientByClientIDFunc: func(clientID uuid.UUID) (domain.Client, error) {
			return &domain.ClientMock{
				IsNotFoundFunc: func() bool {
					return true
				},
			}, nil
		},
	}

	uc := NewAuthorizationUsecase(mockClientRepo, nil, nil, nil, nil)
	_, err := uc.Consent(uuid.New())
	require.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(*errors.UsecaseError).Code)
	assert.Equal(t, "client not found", err.(*errors.UsecaseError).Message)
}
