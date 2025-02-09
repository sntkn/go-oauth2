package usecase

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/domain"
	"github.com/sntkn/go-oauth2/oauth2/domain/domainservice"
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

	uc := NewAuthorizationUsecase(mockClientRepo, nil, nil)
	_, err := uc.Consent(uuid.New())
	require.NoError(t, err)
}

func TestConsent_FindClientError(t *testing.T) {
	mockClientRepo := &domain.ClientRepositoryMock{
		FindClientByClientIDFunc: func(clientID uuid.UUID) (domain.Client, error) {
			return nil, errors.New("FindClientByClientID error")
		},
	}

	uc := NewAuthorizationUsecase(mockClientRepo, nil, nil)
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

	uc := NewAuthorizationUsecase(mockClientRepo, nil, nil)
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

	uc := NewAuthorizationUsecase(nil, mockCodeRepo, nil)
	_, err := uc.GenerateAuthorizationCode(GenerateAuthorizationCodeParams{})
	require.NoError(t, err)
}

func TestGenerateAuthorizationCode_StoreAuthorizationCodeError(t *testing.T) {
	mockCodeRepo := &domain.AuthorizationCodeRepositoryMock{
		StoreAuthorizationCodeFunc: func(storeAuthorizationCodeParams domain.StoreAuthorizationCodeParams) (string, error) {
			return "", errors.New("StoreAuthorizationCode error")
		},
	}

	uc := NewAuthorizationUsecase(nil, mockCodeRepo, nil)
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

	uc := NewAuthorizationUsecase(nil, mockCodeRepo, nil)
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

	uc := NewAuthorizationUsecase(nil, mockCodeRepo, nil)
	_, err := uc.GenerateAuthorizationCode(GenerateAuthorizationCodeParams{})
	require.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(*errors.UsecaseError).Code)
}

func TestGenerateTokenByCode_Success(t *testing.T) {
	mockCodeRepo := &domain.AuthorizationCodeRepositoryMock{
		FindAuthorizationCodeFunc: func(s string) (domain.AuthorizationCode, error) {
			return &domain.AuthorizationCodeMock{
				IsExpiredFunc: func(t time.Time) bool {
					return false
				},
				GetClientIDFunc: func() uuid.UUID {
					return uuid.New()
				},
				GetUserIDFunc: func() uuid.UUID {
					return uuid.New()
				},
				GetScopeFunc: func() string {
					return "scope"
				},
			}, nil
		},
		RevokeCodeFunc: func(code string) error {
			return nil
		},
	}

	mockTokenService := &domainservice.TokenServiceMock{
		StoreNewTokenFunc: func(clientID, UserID uuid.UUID, scope string) (domain.Token, error) {
			return &domain.TokenMock{
				GetAccessTokenFunc: func() string {
					return "access_token"
				},
			}, nil
		},
		StoreNewRefreshTokenFunc: func(accessToken string) (domain.RefreshToken, error) {
			return &domain.RefreshTokenMock{
				GetRefreshTokenFunc: func() string {
					return "refresh_token"
				},
			}, nil
		},
	}

	uc := NewAuthorizationUsecase(nil, mockCodeRepo, mockTokenService)
	token, rtoken, err := uc.GenerateTokenByCode("code")
	require.NoError(t, err)
	assert.Equal(t, "access_token", token.GetAccessToken())
	assert.Equal(t, "refresh_token", rtoken.GetRefreshToken())
}

func TestGenerateTokenByCode_FindValidAuthorizationCodeError(t *testing.T) {
	mockCodeRepo := &domain.AuthorizationCodeRepositoryMock{
		FindAuthorizationCodeFunc: func(s string) (domain.AuthorizationCode, error) {
			return nil, errors.New("FindValidAuthorizationCode error")
		},
	}

	mockTokenService := &domainservice.TokenServiceMock{}

	uc := NewAuthorizationUsecase(nil, mockCodeRepo, mockTokenService)
	_, _, err := uc.GenerateTokenByCode("code")
	require.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*errors.UsecaseError).Code)
	assert.Equal(t, "FindValidAuthorizationCode error", err.(*errors.UsecaseError).Message)
}

func TestGenerateTokenByCode_CodeIsNil(t *testing.T) {
	mockCodeRepo := &domain.AuthorizationCodeRepositoryMock{
		FindAuthorizationCodeFunc: func(s string) (domain.AuthorizationCode, error) {
			return nil, nil
		},
	}

	mockTokenService := &domainservice.TokenServiceMock{}

	uc := NewAuthorizationUsecase(nil, mockCodeRepo, mockTokenService)
	_, _, err := uc.GenerateTokenByCode("code")
	require.Error(t, err)
	assert.Equal(t, http.StatusForbidden, err.(*errors.UsecaseError).Code)
	assert.Equal(t, "code not found", err.(*errors.UsecaseError).Message)
}

func TestGenerateTokenByCode_CodeIsExpired(t *testing.T) {
	mockCodeRepo := &domain.AuthorizationCodeRepositoryMock{
		FindAuthorizationCodeFunc: func(s string) (domain.AuthorizationCode, error) {
			return &domain.AuthorizationCodeMock{
				IsExpiredFunc: func(t time.Time) bool {
					return true
				},
			}, nil
		},
	}

	mockTokenService := &domainservice.TokenServiceMock{}

	uc := NewAuthorizationUsecase(nil, mockCodeRepo, mockTokenService)
	_, _, err := uc.GenerateTokenByCode("code")
	require.Error(t, err)
	assert.Equal(t, http.StatusForbidden, err.(*errors.UsecaseError).Code)
	assert.Equal(t, "code has expired", err.(*errors.UsecaseError).Message)
}

func TestGenerateTokenByCode_StoreTokenError(t *testing.T) {
	mockCodeRepo := &domain.AuthorizationCodeRepositoryMock{
		FindAuthorizationCodeFunc: func(s string) (domain.AuthorizationCode, error) {
			return &domain.AuthorizationCodeMock{
				IsExpiredFunc: func(t time.Time) bool {
					return false
				},
				GetClientIDFunc: func() uuid.UUID {
					return uuid.New()
				},
				GetUserIDFunc: func() uuid.UUID {
					return uuid.New()
				},
				GetScopeFunc: func() string {
					return "scope"
				},
			}, nil
		},
		RevokeCodeFunc: func(code string) error {
			return nil
		},
	}

	mockTokenService := &domainservice.TokenServiceMock{
		StoreNewTokenFunc: func(clientID, UserID uuid.UUID, scope string) (domain.Token, error) {
			return nil, errors.New("StoreNewToken error")
		},
	}

	uc := NewAuthorizationUsecase(nil, mockCodeRepo, mockTokenService)
	_, _, err := uc.GenerateTokenByCode("code")
	require.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*errors.UsecaseError).Code)
	assert.Equal(t, "StoreNewToken error", err.(*errors.UsecaseError).Message)
}

func TestGenerateTokenByCode_StoreRefreshTokenError(t *testing.T) {
	mockCodeRepo := &domain.AuthorizationCodeRepositoryMock{
		FindAuthorizationCodeFunc: func(s string) (domain.AuthorizationCode, error) {
			return &domain.AuthorizationCodeMock{
				IsExpiredFunc: func(t time.Time) bool {
					return false
				},
				GetClientIDFunc: func() uuid.UUID {
					return uuid.New()
				},
				GetUserIDFunc: func() uuid.UUID {
					return uuid.New()
				},
				GetScopeFunc: func() string {
					return "scope"
				},
			}, nil
		},
	}

	mockTokenService := &domainservice.TokenServiceMock{
		StoreNewTokenFunc: func(clientID, UserID uuid.UUID, scope string) (domain.Token, error) {
			return &domain.TokenMock{
				GetAccessTokenFunc: func() string {
					return "access_token"
				},
			}, nil
		},
		StoreNewRefreshTokenFunc: func(accessToken string) (domain.RefreshToken, error) {
			return nil, errors.New("StoreNewRefreshToken error")
		},
	}

	uc := NewAuthorizationUsecase(nil, mockCodeRepo, mockTokenService)
	_, _, err := uc.GenerateTokenByCode("code")
	require.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*errors.UsecaseError).Code)
	assert.Equal(t, "StoreNewRefreshToken error", err.(*errors.UsecaseError).Message)
}

func TestGenerateTokenByCode_RevokeCodeError(t *testing.T) {
	mockCodeRepo := &domain.AuthorizationCodeRepositoryMock{
		FindAuthorizationCodeFunc: func(s string) (domain.AuthorizationCode, error) {
			return &domain.AuthorizationCodeMock{
				IsExpiredFunc: func(t time.Time) bool {
					return false
				},
				GetClientIDFunc: func() uuid.UUID {
					return uuid.New()
				},
				GetUserIDFunc: func() uuid.UUID {
					return uuid.New()
				},
				GetScopeFunc: func() string {
					return "scope"
				},
			}, nil
		},
		RevokeCodeFunc: func(code string) error {
			return errors.New("RevokeCode error")
		},
	}

	mockTokenService := &domainservice.TokenServiceMock{
		StoreNewTokenFunc: func(clientID, UserID uuid.UUID, scope string) (domain.Token, error) {
			return &domain.TokenMock{
				GetAccessTokenFunc: func() string {
					return "access_token"
				},
			}, nil
		},
		StoreNewRefreshTokenFunc: func(accessToken string) (domain.RefreshToken, error) {
			return &domain.RefreshTokenMock{}, nil
		},
	}

	uc := NewAuthorizationUsecase(nil, mockCodeRepo, mockTokenService)
	_, _, err := uc.GenerateTokenByCode("code")
	require.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*errors.UsecaseError).Code)
	assert.Equal(t, "RevokeCode error", err.(*errors.UsecaseError).Message)
}

func TestGenerateTokenByRefreshToken_Success(t *testing.T) {
	mockTokenService := &domainservice.TokenServiceMock{
		FindTokenAndRefreshTokenByRefreshTokenFunc: func(refreshToken string, now time.Time) (domain.Token, domain.RefreshToken, error) {
			return &domain.TokenMock{
					GetAccessTokenFunc: func() string {
						return "access_token"
					},
					GetClientIDFunc: func() uuid.UUID {
						return uuid.New()
					},
					GetUserIDFunc: func() uuid.UUID {
						return uuid.New()
					},
					GetScopeFunc: func() string {
						return "scope"
					},
				}, &domain.RefreshTokenMock{
					GetRefreshTokenFunc: func() string {
						return "refresh_token"
					},
				}, nil
		},
		StoreNewTokenFunc: func(clientID, UserID uuid.UUID, scope string) (domain.Token, error) {
			return &domain.TokenMock{
				GetAccessTokenFunc: func() string {
					return "new_access_token"
				},
			}, nil
		},
		StoreNewRefreshTokenFunc: func(accessToken string) (domain.RefreshToken, error) {
			return &domain.RefreshTokenMock{
				GetRefreshTokenFunc: func() string {
					return "new_refresh_token"
				},
			}, nil
		},
		RevokeTokenFunc: func(accessToken string) error {
			return nil
		},
		RevokeRefreshTokenFunc: func(refreshToken string) error {
			return nil
		},
	}

	uc := NewAuthorizationUsecase(nil, nil, mockTokenService)
	token, rtoken, err := uc.GenerateTokenByRefreshToken("refresh_token")
	require.NoError(t, err)
	assert.Equal(t, "new_access_token", token.GetAccessToken())
	assert.Equal(t, "new_refresh_token", rtoken.GetRefreshToken())
}

func TestGenerateTokenByRefreshToken_FindTokenAndRefreshTokenByRefreshTokenError(t *testing.T) {
	mockTokenService := &domainservice.TokenServiceMock{
		FindTokenAndRefreshTokenByRefreshTokenFunc: func(refreshToken string, now time.Time) (domain.Token, domain.RefreshToken, error) {
			return nil, nil, errors.NewServiceErrorError(errors.ErrCodeInternalServer, "FindTokenAndRefreshTokenByRefreshToken error")
		},
		StoreNewTokenFunc: func(clientID, UserID uuid.UUID, scope string) (domain.Token, error) {
			return &domain.TokenMock{
				GetAccessTokenFunc: func() string {
					return "new_access_token"
				},
			}, nil
		},
		StoreNewRefreshTokenFunc: func(accessToken string) (domain.RefreshToken, error) {
			return &domain.RefreshTokenMock{
				GetRefreshTokenFunc: func() string {
					return "new_refresh_token"
				},
			}, nil
		},
		RevokeTokenFunc: func(accessToken string) error {
			return nil
		},
		RevokeRefreshTokenFunc: func(refreshToken string) error {
			return nil
		},
	}

	uc := NewAuthorizationUsecase(nil, nil, mockTokenService)
	token, rtoken, err := uc.GenerateTokenByRefreshToken("refresh_token")
	require.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*errors.UsecaseError).Code)
	assert.Contains(t, err.(*errors.UsecaseError).Message, "FindTokenAndRefreshTokenByRefreshToken error")
	assert.Nil(t, token)
	assert.Nil(t, rtoken)
}

func TestGenerateTokenByRefreshToken_StoreNewTokenError(t *testing.T) {
	mockTokenService := &domainservice.TokenServiceMock{
		FindTokenAndRefreshTokenByRefreshTokenFunc: func(refreshToken string, now time.Time) (domain.Token, domain.RefreshToken, error) {
			return &domain.TokenMock{
					GetAccessTokenFunc: func() string {
						return "access_token"
					},
					GetClientIDFunc: func() uuid.UUID {
						return uuid.New()
					},
					GetUserIDFunc: func() uuid.UUID {
						return uuid.New()
					},
					GetScopeFunc: func() string {
						return "scope"
					},
				}, &domain.RefreshTokenMock{
					GetRefreshTokenFunc: func() string {
						return "refresh_token"
					},
				}, nil
		},
		StoreNewTokenFunc: func(clientID, UserID uuid.UUID, scope string) (domain.Token, error) {
			return nil, errors.NewServiceErrorError(errors.ErrCodeInternalServer, "StoreNewToken error")
		},
		StoreNewRefreshTokenFunc: func(accessToken string) (domain.RefreshToken, error) {
			return &domain.RefreshTokenMock{
				GetRefreshTokenFunc: func() string {
					return "new_refresh_token"
				},
			}, nil
		},
		RevokeTokenFunc: func(accessToken string) error {
			return nil
		},
		RevokeRefreshTokenFunc: func(refreshToken string) error {
			return nil
		},
	}

	uc := NewAuthorizationUsecase(nil, nil, mockTokenService)
	token, rtoken, err := uc.GenerateTokenByRefreshToken("refresh_token")
	require.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*errors.UsecaseError).Code)
	assert.Contains(t, err.(*errors.UsecaseError).Message, "StoreNewToken error")
	assert.Nil(t, token)
	assert.Nil(t, rtoken)
}

func TestGenerateTokenByRefreshToken_StoreNewRefreshTokenError(t *testing.T) {
	mockTokenService := &domainservice.TokenServiceMock{
		FindTokenAndRefreshTokenByRefreshTokenFunc: func(refreshToken string, now time.Time) (domain.Token, domain.RefreshToken, error) {
			return &domain.TokenMock{
					GetAccessTokenFunc: func() string {
						return "access_token"
					},
					GetClientIDFunc: func() uuid.UUID {
						return uuid.New()
					},
					GetUserIDFunc: func() uuid.UUID {
						return uuid.New()
					},
					GetScopeFunc: func() string {
						return "scope"
					},
				}, &domain.RefreshTokenMock{
					GetRefreshTokenFunc: func() string {
						return "refresh_token"
					},
				}, nil
		},
		StoreNewTokenFunc: func(clientID, UserID uuid.UUID, scope string) (domain.Token, error) {
			return &domain.TokenMock{
				GetAccessTokenFunc: func() string {
					return "new_access_token"
				},
			}, nil
		},
		StoreNewRefreshTokenFunc: func(accessToken string) (domain.RefreshToken, error) {
			return nil, errors.NewServiceErrorError(errors.ErrCodeInternalServer, "StoreNewRefreshToken error")
		},
		RevokeTokenFunc: func(accessToken string) error {
			return nil
		},
		RevokeRefreshTokenFunc: func(refreshToken string) error {
			return nil
		},
	}

	uc := NewAuthorizationUsecase(nil, nil, mockTokenService)
	token, rtoken, err := uc.GenerateTokenByRefreshToken("refresh_token")
	require.Error(t, err)
	assert.Contains(t, err.(*errors.UsecaseError).Message, "StoreNewRefreshToken error")
	assert.Nil(t, token)
	assert.Nil(t, rtoken)
}

func TestGenerateTokenByRefreshToken_RevokeTokenError(t *testing.T) {
	mockTokenService := &domainservice.TokenServiceMock{
		FindTokenAndRefreshTokenByRefreshTokenFunc: func(refreshToken string, now time.Time) (domain.Token, domain.RefreshToken, error) {
			return &domain.TokenMock{
					GetAccessTokenFunc: func() string {
						return "access_token"
					},
					GetClientIDFunc: func() uuid.UUID {
						return uuid.New()
					},
					GetUserIDFunc: func() uuid.UUID {
						return uuid.New()
					},
					GetScopeFunc: func() string {
						return "scope"
					},
				}, &domain.RefreshTokenMock{
					GetRefreshTokenFunc: func() string {
						return "refresh_token"
					},
				}, nil
		},
		StoreNewTokenFunc: func(clientID, UserID uuid.UUID, scope string) (domain.Token, error) {
			return &domain.TokenMock{
				GetAccessTokenFunc: func() string {
					return "new_access_token"
				},
			}, nil
		},
		StoreNewRefreshTokenFunc: func(accessToken string) (domain.RefreshToken, error) {
			return &domain.RefreshTokenMock{
				GetRefreshTokenFunc: func() string {
					return "new_refresh_token"
				},
			}, nil
		},
		RevokeTokenFunc: func(accessToken string) error {
			return errors.NewServiceErrorError(errors.ErrCodeInternalServer, "RevokeToken error")
		},
		RevokeRefreshTokenFunc: func(refreshToken string) error {
			return nil
		},
	}

	uc := NewAuthorizationUsecase(nil, nil, mockTokenService)
	token, rtoken, err := uc.GenerateTokenByRefreshToken("refresh_token")
	require.Error(t, err)
	assert.Contains(t, err.(*errors.UsecaseError).Message, "RevokeToken error")
	assert.Nil(t, token)
	assert.Nil(t, rtoken)
}

func TestGenerateTokenByRefreshToken_RevokeRefreshTokenError(t *testing.T) {
	mockTokenService := &domainservice.TokenServiceMock{
		FindTokenAndRefreshTokenByRefreshTokenFunc: func(refreshToken string, now time.Time) (domain.Token, domain.RefreshToken, error) {
			return &domain.TokenMock{
					GetAccessTokenFunc: func() string {
						return "access_token"
					},
					GetClientIDFunc: func() uuid.UUID {
						return uuid.New()
					},
					GetUserIDFunc: func() uuid.UUID {
						return uuid.New()
					},
					GetScopeFunc: func() string {
						return "scope"
					},
				}, &domain.RefreshTokenMock{
					GetRefreshTokenFunc: func() string {
						return "refresh_token"
					},
				}, nil
		},
		StoreNewTokenFunc: func(clientID, UserID uuid.UUID, scope string) (domain.Token, error) {
			return &domain.TokenMock{
				GetAccessTokenFunc: func() string {
					return "new_access_token"
				},
			}, nil
		},
		StoreNewRefreshTokenFunc: func(accessToken string) (domain.RefreshToken, error) {
			return &domain.RefreshTokenMock{
				GetRefreshTokenFunc: func() string {
					return "new_refresh_token"
				},
			}, nil
		},
		RevokeTokenFunc: func(accessToken string) error {
			return errors.NewServiceErrorError(errors.ErrCodeInternalServer, "RevokeRefreshToken error")
		},
		RevokeRefreshTokenFunc: func(refreshToken string) error {
			return nil
		},
	}

	uc := NewAuthorizationUsecase(nil, nil, mockTokenService)
	token, rtoken, err := uc.GenerateTokenByRefreshToken("refresh_token")
	require.Error(t, err)
	assert.Contains(t, err.(*errors.UsecaseError).Message, "RevokeRefreshToken error")
	assert.Nil(t, token)
	assert.Nil(t, rtoken)
}
