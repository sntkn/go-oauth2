package usecases

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/pkg/str"
	"golang.org/x/crypto/bcrypt"
)

type AuthorizationInput struct {
	Email       string
	Password    string
	Scope       string
	RedirectURI string
	ClientID    string
	Expires     int
}

type Authorization struct {
	cfg *config.Config
	db  repository.OAuth2Repository
}

func NewAuthorization(cfg *config.Config, db repository.OAuth2Repository) *Authorization {
	return &Authorization{
		cfg: cfg,
		db:  db,
	}
}

func (u *Authorization) Invoke(input AuthorizationInput) (string, error) {
	// validate user credentials
	user, err := u.db.FindUserByEmail(input.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.NewUsecaseError(http.StatusBadRequest, err.Error())
		}
		return "", errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	// パスワードを比較して認証
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return "", errors.NewUsecaseError(http.StatusBadRequest, err.Error())
	}

	// create code
	expired := time.Now().Add(time.Duration(input.Expires) * time.Second)
	randomStringLen := 32
	randomString, err := str.GenerateRandomString(randomStringLen)
	if err != nil {
		return "", errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	clientID, err := uuid.Parse(input.ClientID)
	if err != nil {
		return "", errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	err = u.db.RegisterOAuth2Code(&repository.Code{
		Code:        randomString,
		ClientID:    clientID,
		UserID:      user.ID,
		Scope:       input.Scope,
		RedirectURI: input.RedirectURI,
		ExpiresAt:   expired,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		return "", errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	return fmt.Sprintf("%s?code=%s", input.RedirectURI, randomString), nil
}
