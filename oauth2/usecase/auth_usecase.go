package usecase

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/domain/model"
	"github.com/sntkn/go-oauth2/oauth2/infrastructure/repository"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

func NewAuthUsecase(userRepo repository.IUserRepository, clientRepo repository.IClientRepository) IAuthUsecase {
	return &AuthUsecase{
		UserRepo:   userRepo,
		ClientRepo: clientRepo,
	}
}

type IAuthUsecase interface {
	AuthenticateUser(username, password string) (*model.User, error)
	AuthenticateClient(clientID uuid.UUID, redirectURI string) (*model.Client, error)
}

type AuthUsecase struct {
	UserRepo   repository.IUserRepository
	ClientRepo repository.IClientRepository
}

func (uc *AuthUsecase) AuthenticateUser(username, password string) (*model.User, error) {
	return uc.UserRepo.FindUserByEmail(username)
}

func (uc *AuthUsecase) AuthenticateClient(clientID uuid.UUID, redirectURI string) (*model.Client, error) {
	client, err := uc.ClientRepo.FindClientByClientID(clientID)
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

	return model.NewClient(client.ID, client.Name, client.RedirectURIs, client.CreatedAt, client.UpdatedAt), nil
}
