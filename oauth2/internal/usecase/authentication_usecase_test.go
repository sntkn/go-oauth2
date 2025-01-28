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

func TestAuthenticateClient_Success(t *testing.T) {
	mockUserRepo := &domain.UserRepositoryMock{}

	mockClientRepo := &domain.ClientRepositoryMock{
		FindClientByClientIDFunc: func(clientID uuid.UUID) (domain.Client, error) {
			return &domain.ClientMock{
				IsNotFoundFunc: func() bool {
					return false
				},
				IsRedirectURIMatchFunc: func(uri string) bool {
					return true
				},
			}, nil
		},
	}

	uc := NewAuthenticationUsecase(mockUserRepo, mockClientRepo)
	_, err := uc.AuthenticateClient(uuid.New(), "https://example.com/callback")
	require.NoError(t, err)
}

func TestAuthenticateClient_ClientNotFound(t *testing.T) {
	mockUserRepo := &domain.UserRepositoryMock{}

	mockClientRepo := &domain.ClientRepositoryMock{
		FindClientByClientIDFunc: func(clientID uuid.UUID) (domain.Client, error) {
			return &domain.ClientMock{
				IsNotFoundFunc: func() bool {
					return true
				},
				IsRedirectURIMatchFunc: func(uri string) bool {
					return true
				},
			}, nil
		},
	}

	uc := NewAuthenticationUsecase(mockUserRepo, mockClientRepo)
	_, err := uc.AuthenticateClient(uuid.New(), "https://example.com/callback")
	require.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(*errors.UsecaseError).Code)
	assert.Equal(t, "client not found", err.(*errors.UsecaseError).Message)
}

func TestAuthenticateClient_RedirectURINotMatch(t *testing.T) {
	mockUserRepo := &domain.UserRepositoryMock{}

	mockClientRepo := &domain.ClientRepositoryMock{
		FindClientByClientIDFunc: func(clientID uuid.UUID) (domain.Client, error) {
			return &domain.ClientMock{
				IsNotFoundFunc: func() bool {
					return false
				},
				IsRedirectURIMatchFunc: func(uri string) bool {
					return false
				},
			}, nil
		},
	}

	uc := NewAuthenticationUsecase(mockUserRepo, mockClientRepo)
	_, err := uc.AuthenticateClient(uuid.New(), "https://example.com/callback")
	require.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(*errors.UsecaseError).Code)
	assert.Equal(t, "redirect uri does not match", err.(*errors.UsecaseError).Message)
}

func TestAuthenticateClient_FindClientError(t *testing.T) {
	mockUserRepo := &domain.UserRepositoryMock{}

	mockClientRepo := &domain.ClientRepositoryMock{
		FindClientByClientIDFunc: func(clientID uuid.UUID) (domain.Client, error) {
			return nil, errors.New("FindClientByClientID error")
		},
	}

	uc := NewAuthenticationUsecase(mockUserRepo, mockClientRepo)
	_, err := uc.AuthenticateClient(uuid.New(), "https://example.com/callback")
	require.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*errors.UsecaseError).Code)
	assert.Equal(t, "FindClientByClientID error", err.(*errors.UsecaseError).Message)
}

func TestAuthenticateUser_Success(t *testing.T) {
	mockUserRepo := &domain.UserRepositoryMock{
		FindUserByEmailFunc: func(email string) (domain.User, error) {
			return &domain.UserMock{
				IsNotFoundFunc: func() bool {
					return false
				},
				IsPasswordMatchFunc: func(password string) bool {
					return true
				},
			}, nil
		},
	}

	mockClientRepo := &domain.ClientRepositoryMock{}

	uc := NewAuthenticationUsecase(mockUserRepo, mockClientRepo)
	_, err := uc.AuthenticateUser("test@example.com", "password123")
	require.NoError(t, err)
}

func TestAuthenticateUser_UserNotFound(t *testing.T) {
	mockUserRepo := &domain.UserRepositoryMock{
		FindUserByEmailFunc: func(email string) (domain.User, error) {
			return &domain.UserMock{
				IsNotFoundFunc: func() bool {
					return true
				},
			}, nil
		},
	}

	mockClientRepo := &domain.ClientRepositoryMock{}

	uc := NewAuthenticationUsecase(mockUserRepo, mockClientRepo)
	result, err := uc.AuthenticateUser("test@example.com", "password123")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, http.StatusFound, err.(*errors.UsecaseError).Code)
	assert.Equal(t, "user or password not match", err.(*errors.UsecaseError).Message)
	assert.Equal(t, "client/signin", err.(*errors.UsecaseError).RedirectURI)
}

func TestAuthenticateUser_PasswordMismatch(t *testing.T) {
	mockUserRepo := &domain.UserRepositoryMock{
		FindUserByEmailFunc: func(email string) (domain.User, error) {
			return &domain.UserMock{
				IsNotFoundFunc: func() bool {
					return false
				},
				IsPasswordMatchFunc: func(password string) bool {
					return false
				},
			}, nil
		},
	}

	mockClientRepo := &domain.ClientRepositoryMock{}

	uc := NewAuthenticationUsecase(mockUserRepo, mockClientRepo)
	result, err := uc.AuthenticateUser("test@example.com", "password123")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, http.StatusFound, err.(*errors.UsecaseError).Code)
	assert.Equal(t, "user or password not match", err.(*errors.UsecaseError).Message)
	assert.Equal(t, "client/signin", err.(*errors.UsecaseError).RedirectURI)
}

func TestAuthenticateUser_FindUserError(t *testing.T) {
	mockUserRepo := &domain.UserRepositoryMock{
		FindUserByEmailFunc: func(email string) (domain.User, error) {
			return nil, errors.New("FindUserByEmail error")
		},
	}

	mockClientRepo := &domain.ClientRepositoryMock{}

	uc := NewAuthenticationUsecase(mockUserRepo, mockClientRepo)
	result, err := uc.AuthenticateUser("test@example.com", "password123")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, http.StatusInternalServerError, err.(*errors.UsecaseError).Code)
	assert.Equal(t, "FindUserByEmail error", err.(*errors.UsecaseError).Message)
}
