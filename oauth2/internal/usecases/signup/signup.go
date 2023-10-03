package signup

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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
	form, err := GetFormData(c, s)
	if err != nil {
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusBadRequest, "500.html", gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "signup.html", gin.H{"f": form})
}

func GetFormData(c *gin.Context, s *session.Session) (*RegistrationForm, error) {
	var d RegistrationForm
	sessData, err := s.GetSessionData(c, "create_user_form")
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
