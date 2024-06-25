package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/pkg/redis"
)

type AuthorizationInput struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

type SigninForm struct {
	Email string `form:"email"`
	Error string
}

func AuthrozationHandler(redisCli *redis.RedisCli, db *repository.Repository, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input AuthorizationInput

		if err := c.ShouldBind(&input); err != nil {
			c.Redirect(http.StatusFound, "/signin")
			return
		}

		s := session.NewSession(c, redisCli, cfg.SessionExpires)

		if err := s.SetNamedSessionData(c, "signin_form", SigninForm{
			Email: input.Email,
		}); err != nil {
			c.Error(errors.WithStack(err))
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
			return
		}

		redirectURI, err := usecases.NewAuthorization(redisCli, db, cfg, s).Invoke(c, input.Email, input.Password)
		if err != nil {
			if usecaseErr, ok := err.(*cerrs.UsecaseError); ok {
				switch usecaseErr.Code {
				case http.StatusFound:
					c.Redirect(http.StatusFound, "/signin")

				case http.StatusInternalServerError:
					c.Error(errors.WithStack(err)) // TODO: trigger usecase
					c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": usecaseErr.Error()})
				}
			} else {
				c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
			}
			return
		}

		c.Redirect(http.StatusFound, redirectURI)
	}
}
