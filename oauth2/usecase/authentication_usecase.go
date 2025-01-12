package usecase

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/domain/authentication"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func NewAuthenticationUsecase(repo authentication.IAuthenticationRepository) IAuthenticationUsecase {
	return &AuthenticationUsecase{
		repo: repo,
	}
}

type IAuthenticationUsecase interface {
	AuthenticateUser(email, password string) (*authentication.User, error)
	AuthenticateClient(clientID uuid.UUID, redirectURI string) (*authentication.Client, error)
}

type AuthenticationUsecase struct {
	repo authentication.IAuthenticationRepository
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

func (uc *AuthenticationUsecase) AuthenticateUser(email, password string) (*authentication.User, error) {
	// validate user credentials
	user, err := uc.repo.FindUserByEmail(email)

	if err != nil {
		return nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	// ユーザーが存在しない場合はエラー
	if user.ID == uuid.Nil {
		return nil, errors.NewUsecaseErrorWithRedirectURI(http.StatusFound, "user not found", "client/signin")
	}

	// パスワードを比較して認証
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.NewUsecaseErrorWithRedirectURI(http.StatusFound, err.Error(), "client/signin")
	}

	return authentication.NewUser(user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt), nil
}
