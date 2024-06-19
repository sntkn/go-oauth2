package usecases

import (
	"fmt"
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/redis"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
)

type AuthorizeInput struct {
	ResponseType string `form:"response_type"`
	ClientID     string `form:"client_id"`
	Scope        string `form:"scope"`
	RedirectURI  string `form:"redirect_uri"`
	State        string `form:"state"`
}

type SigninForm struct {
	Email string `form:"email"`
	Error string
}

type Signin struct {
	redisCli *redis.RedisCli
}

func NewSignin(redisCli *redis.RedisCli) *Signin {
	return &Signin{
		redisCli: redisCli,
	}
}

func (u *Signin) Invoke(c *gin.Context) {
	s := session.NewSession(c, u.redisCli)

	var input AuthorizeInput
	if err := s.GetNamedSessionData(c, "auth", &input); err != nil {
		c.Error(err)
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		return
	}

	if input.ClientID == "" {
		err := fmt.Errorf("invalid client_id")
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		return
	}

	var form SigninForm
	if err := s.FlushNamedSessionData(c, "signin_form", &form); err != nil {
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "signin.html", gin.H{"f": form})
}
