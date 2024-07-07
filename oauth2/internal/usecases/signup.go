package usecases

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/entity"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

type Signup struct {
	sess *session.Session
	cfg  *config.Config
}

func NewSignup(cfg *config.Config, sess *session.Session) *Signup {
	return &Signup{
		sess: sess,
		cfg:  cfg,
	}
}

func (u *Signup) Invoke(c *gin.Context) (entity.SessionRegistrationForm, error) {
	var form entity.SessionRegistrationForm
	if err := u.sess.FlushNamedSessionData(c, "signup_form", &form); err != nil {
		return form, cerrs.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}
	return form, nil
}
