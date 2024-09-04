package usecases

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/pkg/str"
	"golang.org/x/crypto/bcrypt"
)

type Authorization struct {
	cfg  *config.Config
	db   repository.OAuth2Repository
	sess session.SessionClient
}

func NewAuthorization(cfg *config.Config, db repository.OAuth2Repository, sess session.SessionClient) *Authorization {
	return &Authorization{
		cfg:  cfg,
		db:   db,
		sess: sess,
	}
}

func (u *Authorization) Invoke(c *gin.Context, email, password string) (string, error) {
	// validate user credentials
	user, err := u.db.FindUserByEmail(email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.NewUsecaseError(http.StatusBadRequest, err.Error())
		}
		return "", errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	// パスワードを比較して認証
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.NewUsecaseError(http.StatusBadRequest, err.Error())
	}

	var d AuthorizeInput
	if err = u.sess.GetNamedSessionData(c, "auth", &d); err != nil {
		return "", errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	// create code
	expired := time.Now().Add(time.Duration(u.cfg.AuthCodeExpires) * time.Second)
	randomStringLen := 32
	randomString, err := str.GenerateRandomString(randomStringLen)
	if err != nil {
		return "", errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	clientID, err := uuid.Parse(d.ClientID)
	if err != nil {
		return "", errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	err = u.db.RegisterOAuth2Code(&repository.Code{
		Code:        randomString,
		ClientID:    clientID,
		UserID:      user.ID,
		Scope:       d.Scope,
		RedirectURI: d.RedirectURI,
		ExpiresAt:   expired,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		return "", errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	if err := u.sess.DelSessionData(c, "auth"); err != nil {
		return "", errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	return fmt.Sprintf("%s?code=%s", d.RedirectURI, randomString), nil
}
