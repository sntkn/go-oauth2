package usecases

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

type RegistrationData struct {
	Name  string
	Email string
}

type CreateUser struct {
	cfg  *config.Config
	db   repository.OAuth2Repository
	sess session.SessionClient
}

func NewCreateUser(cfg *config.Config, db repository.OAuth2Repository, sess session.SessionClient) *CreateUser {
	return &CreateUser{
		cfg:  cfg,
		db:   db,
		sess: sess,
	}
}

func (u *CreateUser) Invoke(c *gin.Context, user repository.User) error {

	if err := u.sess.SetNamedSessionData(c, "signup_form", RegistrationData{
		Name:  user.Name,
		Email: user.Email,
	}); err != nil {
		return cerrs.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	// check email is existing
	eu, err := u.db.ExistsUserByEmail(user.Email)
	if err != nil {
		return cerrs.NewUsecaseError(http.StatusInternalServerError, err.Error())
	} else if eu {
		return cerrs.NewUsecaseError(http.StatusFound, "input email already exists")
	}

	if err := u.db.CreateUser(&user); err != nil {
		return cerrs.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	if err := u.sess.DelSessionData(c, "signup_form"); err != nil {
		return cerrs.NewUsecaseError(http.StatusFound, "input email already exists")
	}

	return nil
}
