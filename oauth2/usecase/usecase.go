package usecase

import "github.com/sntkn/go-oauth2/oauth2/domain/model"

type TokenUsecase interface {
	GenerateAccessToken(user *model.User, client *model.Client, scopes []string) (*model.AccessToken, *model.RefreshToken, error)
	ValidateAccessToken(token string) (*model.AccessToken, error)
}
