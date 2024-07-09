package usecases

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/entity"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/pkg/redis"
)

type SigninForm struct {
	Email string `form:"email"`
	Error string
}

type Signin struct {
	redisCli *redis.RedisCli
	cfg      *config.Config
}

func NewSignin(redisCli *redis.RedisCli, cfg *config.Config) *Signin {
	return &Signin{
		redisCli: redisCli,
		cfg:      cfg,
	}
}

func (u *Signin) Invoke(c *gin.Context) (entity.SessionSigninForm, error) {
	s := session.NewSession(c, u.redisCli, u.cfg.SessionExpires)

	var form entity.SessionSigninForm
	var input AuthorizeInput

	if err := s.GetNamedSessionData(c, "auth", &input); err != nil {
		return form, cerrs.NewUsecaseError(http.StatusBadRequest, err.Error())
	}

	if input.ClientID == "" {
		return form, cerrs.NewUsecaseError(http.StatusBadRequest, "invalid client_id")
	}

	if err := s.FlushNamedSessionData(c, "signin_form", &form); err != nil {
		return form, cerrs.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}
	return form, nil
}
