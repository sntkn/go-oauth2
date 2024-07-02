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

type Signup struct {
	redisCli *redis.RedisCli
	cfg      *config.Config
}

func NewSignup(redisCli *redis.RedisCli, cfg *config.Config) *Signup {
	return &Signup{
		redisCli: redisCli,
		cfg:      cfg,
	}
}

func (u *Signup) Invoke(c *gin.Context) (entity.SessionRegistrationForm, error) {
	s := session.NewSession(c, u.redisCli, u.cfg.SessionExpires)
	var form entity.SessionRegistrationForm
	if err := s.FlushNamedSessionData(c, "signup_form", &form); err != nil {
		return form, cerrs.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}
	return form, nil
}
