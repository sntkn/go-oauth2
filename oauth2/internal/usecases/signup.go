package usecases

import (
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/redis"
)

type RegistrationForm struct {
	Name  string `form:"name"`
	Email string `form:"email"`
	Error string
}

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

func (u *Signup) Invoke(c *gin.Context) {
	s := session.NewSession(c, u.redisCli, u.cfg.SessionExpires)
	var form RegistrationForm
	if err := s.FlushNamedSessionData(c, "signup_form", &form); err != nil {
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusBadRequest, "500.html", gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "signup.html", gin.H{"f": form})
}
