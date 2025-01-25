package usecase

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/domain"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

func NewAuthenticationUsecase(userRepo domain.UserRepository, clientRepo domain.ClientRepository) IAuthenticationUsecase {
	return &AuthenticationUsecase{
		userRepo:   userRepo,
		clientRepo: clientRepo,
	}
}

type IAuthenticationUsecase interface {
	AuthenticateUser(email, password string) (*domain.User, error)
	AuthenticateClient(clientID uuid.UUID, redirectURI string) (*domain.Client, error)
}

type AuthenticationUsecase struct {
	userRepo   domain.UserRepository
	clientRepo domain.ClientRepository
}

func (uc *AuthenticationUsecase) AuthenticateClient(clientID uuid.UUID, redirectURI string) (*domain.Client, error) {
	cli, err := uc.clientRepo.FindClientByClientID(clientID)
	if err != nil {
		return nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	client := domain.NewClient(cli.ID, cli.Name, cli.RedirectURIs, cli.CreatedAt, cli.UpdatedAt)

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

func (uc *AuthenticationUsecase) AuthenticateUser(email, password string) (*domain.User, error) {
	// validate user credentials
	u, err := uc.userRepo.FindUserByEmail(email)

	if err != nil {
		return nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	user := domain.NewUser(u.ID, u.Name, u.Email, u.Password, u.CreatedAt, u.UpdatedAt)

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
