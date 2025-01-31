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

func TestGenerateAuthorizationCode_Success(t *testing.T) {
	mockCodeRepo := &domain.AuthorizationCodeRepositoryMock{
		FindAuthorizationCodeFunc: func(s string) (domain.AuthorizationCode, error) {
			return &domain.AuthorizationCodeMock{}, nil
		},
		StoreAuthorizationCodeFunc: func(storeAuthorizationCodeParams domain.StoreAuthorizationCodeParams) (string, error) {
			return "code", nil
		},
	}

	uc := NewAuthorizationUsecase(nil, mockCodeRepo, nil, nil, nil)
	_, err := uc.GenerateAuthorizationCode(GenerateAuthorizationCodeParams{})
	require.NoError(t, err)
}

func TestGenerateAuthorizationCode_StoreAuthorizationCodeError(t *testing.T) {
	mockCodeRepo := &domain.AuthorizationCodeRepositoryMock{
		StoreAuthorizationCodeFunc: func(storeAuthorizationCodeParams domain.StoreAuthorizationCodeParams) (string, error) {
			return "", errors.New("StoreAuthorizationCode error")
		},
	}

	uc := NewAuthorizationUsecase(nil, mockCodeRepo, nil, nil, nil)
	_, err := uc.GenerateAuthorizationCode(GenerateAuthorizationCodeParams{})
	require.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*errors.UsecaseError).Code)
	assert.Equal(t, "StoreAuthorizationCode error", err.(*errors.UsecaseError).Message)
}

func TestGenerateAuthorizationCode_FindAuthorizationCodeError(t *testing.T) {
	mockCodeRepo := &domain.AuthorizationCodeRepositoryMock{
		FindAuthorizationCodeFunc: func(s string) (domain.AuthorizationCode, error) {
			return nil, errors.New("FindAuthorizationCode error")
		},
		StoreAuthorizationCodeFunc: func(storeAuthorizationCodeParams domain.StoreAuthorizationCodeParams) (string, error) {
			return "code", nil
		},
	}

	uc := NewAuthorizationUsecase(nil, mockCodeRepo, nil, nil, nil)
	_, err := uc.GenerateAuthorizationCode(GenerateAuthorizationCodeParams{})
	require.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*errors.UsecaseError).Code)
	assert.Equal(t, "FindAuthorizationCode error", err.(*errors.UsecaseError).Message)
}

func TestGenerateAuthorizationCode_FindAuthorizationCodeNil(t *testing.T) {
	mockCodeRepo := &domain.AuthorizationCodeRepositoryMock{
		FindAuthorizationCodeFunc: func(s string) (domain.AuthorizationCode, error) {
			return nil, nil
		},
		StoreAuthorizationCodeFunc: func(storeAuthorizationCodeParams domain.StoreAuthorizationCodeParams) (string, error) {
			return "code", nil
		},
	}

	uc := NewAuthorizationUsecase(nil, mockCodeRepo, nil, nil, nil)
	_, err := uc.GenerateAuthorizationCode(GenerateAuthorizationCodeParams{})
	require.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(*errors.UsecaseError).Code)
}
