package usecase

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/domain/authentication"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

func NewAuthenticationUsecase(repo authentication.IAuthenticationRepository) IAuthenticationUsecase {
	return &AuthenticationUsecase{
		repo: repo,
	}
}

type IAuthenticationUsecase interface {
	AuthenticateUser(username, password string) (*authentication.User, error)
	AuthenticateClient(clientID uuid.UUID, redirectURI string) (*authentication.Client, error)
}

type AuthenticationUsecase struct {
	repo authentication.IAuthenticationRepository
}

func (uc *AuthenticationUsecase) AuthenticateUser(username, password string) (*authentication.User, error) {
	user, err := uc.repo.FindUserByEmail(username)
	if err != nil {
		return nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}
	return authentication.NewUser(user.ID, user.Email, user.Name, user.CreatedAt, user.UpdatedAt), nil
}

func (uc *AuthenticationUsecase) AuthenticateClient(clientID uuid.UUID, redirectURI string) (*authentication.Client, error) {
	client, err := uc.repo.FindClientByClientID(clientID)
	if err != nil {
		return nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	// クライアントがない場合はエラー
	if client.ID != clientID {
		return nil, errors.NewUsecaseError(http.StatusBadRequest, "client not found")
	}

	// リダイレクトURIが一致しない場合はエラー
	if client.RedirectURIs != redirectURI {
		return nil, errors.NewUsecaseError(http.StatusBadRequest, "redirect uri does not match")
	}

	return authentication.NewClient(client.ID, client.Name, client.RedirectURIs, client.CreatedAt, client.UpdatedAt), nil
}
