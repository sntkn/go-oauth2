package usecases

import (
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/redis"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
)

type RegistrationForm struct {
	Name  string `form:"name"`
	Email string `form:"email"`
	Error string
}

type Signup struct {
	redisCli *redis.RedisCli
}

func NewSignup(redisCli *redis.RedisCli) *Signup {
	return &Signup{
		redisCli: redisCli,
	}
}

func (u *Signup) Invoke(c *gin.Context) {
	s := session.NewSession(c, u.redisCli)
	var form RegistrationForm
	if err := s.FlushNamedSessionData(c, "signup_form", &form); err != nil {
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusBadRequest, "500.html", gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "signup.html", gin.H{"f": form})
}
