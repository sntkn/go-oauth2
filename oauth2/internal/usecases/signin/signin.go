package signin

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/internal/redis"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
)

type SigninForm struct {
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

	ss, err := s.GetSessionData(c, "auth")
	if err != nil {
		c.Error(err)
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		return
	}

	if len(ss) == 0 {
		c.Error(err)
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		return
	}

	form, err := GetFormData(c, s)
	if err != nil {
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusBadRequest, "500.html", gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "signin.html", gin.H{"f": form})
}

func GetFormData(c *gin.Context, s *session.Session) (*SigninForm, error) {
	var d SigninForm
	sessData, err := s.GetSessionData(c, "signin_form")
	if err != nil {
		return nil, err
	} else if sessData == nil {
		return &d, nil
	}

	err = json.Unmarshal(sessData, &d)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &d, nil
}
