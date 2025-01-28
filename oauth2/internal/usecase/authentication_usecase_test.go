package usecase

import (
	"net/http"
	"testing"

	"github.com/sntkn/go-oauth2/oauth2/domain"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	assert.Equal(t, "", err.(*errors.UsecaseError).RedirectURI)
}
