package usecase

import (
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/domain/model"
	"github.com/sntkn/go-oauth2/oauth2/interface/repository"
)

func NewAuthUsecase(userRepo repository.IUserRepository, clientRepo repository.IClientRepository) IAuthUsecase {
	return &AuthUsecase{
		UserRepo:   userRepo,
		ClientRepo: clientRepo,
	}
}

type IAuthUsecase interface {
	AuthenticateUser(username, password string) (*model.User, error)
	AuthenticateClient(clientID uuid.UUID, clientSecret string) (*model.Client, error)
}

type AuthUsecase struct {
	UserRepo   repository.IUserRepository
	ClientRepo repository.IClientRepository
}

func (uc *AuthUsecase) AuthenticateUser(username, password string) (*model.User, error) {
	return uc.UserRepo.FindUserByEmail(username)
}

func (uc *AuthUsecase) AuthenticateClient(clientID uuid.UUID, clientSecret string) (*model.Client, error) {
	return uc.ClientRepo.FindClientByClientID(clientID)
}
