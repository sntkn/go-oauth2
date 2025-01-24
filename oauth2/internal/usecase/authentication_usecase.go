package usecase

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/domain/authentication"
	"github.com/sntkn/go-oauth2/oauth2/domain/client"
	"github.com/sntkn/go-oauth2/oauth2/domain/user"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

func NewAuthenticationUsecase(userRepo user.UserRepository, clientRepo client.ClientRepository) IAuthenticationUsecase {
	return &AuthenticationUsecase{
		userRepo:   userRepo,
		clientRepo: clientRepo,
	}
}

type IAuthenticationUsecase interface {
	AuthenticateUser(email, password string) (*authentication.User, error)
	AuthenticateClient(clientID uuid.UUID, redirectURI string) (*authentication.Client, error)
}

type AuthenticationUsecase struct {
	userRepo   user.UserRepository
	clientRepo client.ClientRepository
}

func (uc *AuthenticationUsecase) AuthenticateClient(clientID uuid.UUID, redirectURI string) (*authentication.Client, error) {
	cli, err := uc.clientRepo.FindClientByClientID(clientID)
	if err != nil {
		return nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	client := authentication.NewClient(cli.ID, cli.Name, cli.RedirectURIs, cli.CreatedAt, cli.UpdatedAt)

	// クライアントがない場合はエラー
	if client.IsNotFound() {
		return nil, errors.NewUsecaseError(http.StatusBadRequest, "client not found")
	}

	// リダイレクトURIが一致しない場合はエラー
	if !client.IsRedirectURIMatch(redirectURI) {
		return nil, errors.NewUsecaseError(http.StatusBadRequest, "redirect uri does not match")
	}

	return client, nil
}

func (uc *AuthenticationUsecase) AuthenticateUser(email, password string) (*authentication.User, error) {
	// validate user credentials
	u, err := uc.userRepo.FindUserByEmail(email)

	if err != nil {
		return nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	user := authentication.NewUser(u.ID, u.Name, u.Email, u.Password, u.CreatedAt, u.UpdatedAt)

	// ユーザーが存在しない場合はエラー
	if user.IsNotFound() {
		return nil, errors.NewUsecaseErrorWithRedirectURI(http.StatusFound, "user or password not match", "client/signin")
	}

	// パスワードを比較して認証
	if !user.IsPasswordMatch(password) {
		return nil, errors.NewUsecaseErrorWithRedirectURI(http.StatusFound, "user or password not match", "client/signin")
	}

	return user, nil
}
