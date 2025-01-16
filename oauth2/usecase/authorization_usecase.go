package usecase

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/domain/authorization"
	"github.com/sntkn/go-oauth2/oauth2/infrastructure/model"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

type IAuthorizationUsecase interface {
	Consent(clientID uuid.UUID) (*authorization.Client, error)
	GenerateAuthorizationCode(GenerateAuthorizationCodeParams) (*authorization.AuthorizationCode, error)
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

func (uc *AuthorizationUsecase) Consent(clientID uuid.UUID) (*authorization.Client, error) {
	cli, err := uc.repo.FindClientByClientID(clientID)
	if err != nil {
		return nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	client := authorization.NewClient(cli.ID, cli.Name, cli.RedirectURIs, cli.CreatedAt, cli.UpdatedAt)
	if client.IsNotFound() {
		return nil, errors.NewUsecaseError(http.StatusBadRequest, "client not found")
	}

	return client, nil
}

type GenerateAuthorizationCodeParams struct {
	UserID      string
	ClientID    string
	RedirectURI string
	Scope       string
	Expires     int
}

func (uc *AuthorizationUsecase) GenerateAuthorizationCode(p GenerateAuthorizationCodeParams) (*authorization.AuthorizationCode, error) {
	randomString, err := authorization.GenerateCode()
	if err != nil {
		return nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	clientID, err := uuid.Parse(p.ClientID)
	if err != nil {
		return nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}
	userID, err := uuid.Parse(p.UserID)
	if err != nil {
		return nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	code := &model.AuthorizationCode{
		Code:        randomString,
		ClientID:    clientID,
		UserID:      userID,
		Scope:       p.Scope,
		RedirectURI: p.RedirectURI,
		ExpiresAt:   time.Now().Add(time.Duration(p.Expires) * time.Second),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = uc.repo.StoreAuthorizationCode(code)
	if err != nil {
		return nil, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	client := authorization.NewAuthorizationCode(
		code.Code,
		code.ClientID,
		code.UserID,
		code.Scope,
		code.RedirectURI,
		code.ExpiresAt,
		code.CreatedAt,
		code.UpdatedAt,
	)

	return client, nil
}
