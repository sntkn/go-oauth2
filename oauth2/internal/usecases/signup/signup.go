package signup

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

type UseCase struct {
	redisCli *redis.RedisCli
}

func NewUseCase(redisCli *redis.RedisCli) *UseCase {
	return &UseCase{
		redisCli: redisCli,
	}
}

func (u *UseCase) Run(c *gin.Context) {
	s := session.NewSession(c, u.redisCli)
	var form RegistrationForm
	if err := s.FlushNamedSessionData(c, "signup_form", &form); err != nil {
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusBadRequest, "500.html", gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "signup.html", gin.H{"f": form})
}
