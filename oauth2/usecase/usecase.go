package usecase

import "github.com/sntkn/go-oauth2/oauth2/domain/model"

type AuthorizationUsecase interface {
	GenerateAuthorizationCode(user *model.User, client *model.Client, scopes []string) (*model.AuthorizationCode, error)
	ValidateAuthorizationCode(code string, clientID string) (*model.AuthorizationCode, error)
}

type TokenUsecase interface {
	GenerateAccessToken(user *model.User, client *model.Client, scopes []string) (*model.AccessToken, *model.RefreshToken, error)
	ValidateAccessToken(token string) (*model.AccessToken, error)
}
