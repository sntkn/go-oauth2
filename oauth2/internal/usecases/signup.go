package usecases

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/entity"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

type Signup struct {
	cfg  *config.Config
	sess session.SessionClient
}

func NewSignup(cfg *config.Config, sess session.SessionClient) *Signup {
	return &Signup{
		cfg:  cfg,
		sess: sess,
	}
}

func (u *Signup) Invoke(c *gin.Context) (entity.SessionRegistrationForm, error) {
	var form entity.SessionRegistrationForm
	if err := u.sess.FlushNamedSessionData(c, "signup_form", &form); err != nil {
		return form, errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}
	return form, nil
}
