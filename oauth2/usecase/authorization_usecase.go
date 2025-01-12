package usecase

import "github.com/sntkn/go-oauth2/oauth2/domain/authorization"

type IAuthorizationUsecase interface {
	// GenerateAuthorizationCode(user *model.User, client *model.Client, scopes []string) (*model.AuthorizationCode, error)
	// ValidateAuthorizationCode(code string, clientID string) (*model.AuthorizationCode, error)
}

func NewAuthorizationUsecase(repo authorization.IAuthorizationRepository) IAuthorizationUsecase {
	return &AuthorizationUsecase{
		repo: repo,
	}
}

type AuthorizationUsecase struct {
	repo authorization.IAuthorizationRepository
}
