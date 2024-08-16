package usecases

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/entity"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

type Signin struct {
	cfg  *config.Config
	sess session.SessionClient
}

func NewSignin(cfg *config.Config, sess session.SessionClient) *Signin {
	return &Signin{
		cfg:  cfg,
		sess: sess,
	}
}

func (u *Signin) Invoke(c *gin.Context) (entity.SessionSigninForm, error) {
	var form entity.SessionSigninForm
	var input AuthorizeInput

	if err := u.sess.GetNamedSessionData(c, "auth", &input); err != nil {
		return form, cerrs.NewUsecaseError(http.StatusBadRequest, err.Error())
	}

	if input.ClientID == "" {
		return form, cerrs.NewUsecaseError(http.StatusBadRequest, "invalid client_id")
	}

	if err := u.sess.FlushNamedSessionData(c, "signin_form", &form); err != nil {
		return form, cerrs.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}
	return form, nil
}
