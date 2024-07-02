package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/pkg/redis"
)

type RegistrationForm struct {
	Name  string `form:"name"`
	Email string `form:"email"`
	Error string
}

func SignupHandler(redisCli *redis.RedisCli, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		form, err := usecases.NewSignup(redisCli, cfg).Invoke(c)
		if err != nil {
			if usecaseErr, ok := err.(*cerrs.UsecaseError); ok {
				switch usecaseErr.Code {
				case http.StatusInternalServerError:
					c.Error(errors.WithStack(err)) // TODO: trigger usecase
					c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": usecaseErr.Error()})
				}
				return
			}
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
			return
		}
		c.HTML(http.StatusOK, "signup.html", gin.H{"f": form})
	}
}
