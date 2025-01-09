package usecase

import "github.com/sntkn/go-oauth2/oauth2/domain/model"

type AuthUsecase interface {
	AuthenticateUser(username, password string) (*model.User, error)
	AuthenticateClient(clientID, clientSecret string) (*model.Client, error)
}

type AuthorizationUsecase interface {
	GenerateAuthorizationCode(user *model.User, client *model.Client, scopes []string) (*model.AuthorizationCode, error)
	ValidateAuthorizationCode(code string, clientID string) (*model.AuthorizationCode, error)
}

type TokenUsecase interface {
	GenerateAccessToken(user *model.User, client *model.Client, scopes []string) (*model.AccessToken, *model.RefreshToken, error)
	ValidateAccessToken(token string) (*model.AccessToken, error)
}
