package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

type RegistrationForm struct {
	Name  string `form:"name"`
	Email string `form:"email"`
	Error string
}

func SignupHandler(sessionCreator session.Creator, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := sessionCreator(c)
		form, err := usecases.NewSignup(cfg, s).Invoke(c)
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

		mess, err := s.PullSessionData(c, "flushMessage")
		if err != nil {
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		}

		c.HTML(http.StatusOK, "signup.html", gin.H{"f": form, "m": string(mess)})
	}
}
