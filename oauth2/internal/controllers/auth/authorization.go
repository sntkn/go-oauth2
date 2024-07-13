package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/internal/flashmessage"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

type AuthorizationInput struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

type SigninForm struct {
	Email string `form:"email"`
	Error string
}

func AuthrozationHandler(db *repository.Repository, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		s, err := session.GetSession(c)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
			return
		}
		//mess, err := flashmessage.GetMessage(c)
		//if err != nil {
		//	c.Error(errors.WithStack(err)) // TODO: trigger usecase
		//	c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		//	return
		//}

		var input AuthorizationInput

		if err := s.SetNamedSessionData(c, "signin_form", SigninForm{
			Email: input.Email,
		}); err != nil {
			c.Error(errors.WithStack(err))
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
			return
		}

		if err := c.ShouldBind(&input); err != nil {
			if err := flashmessage.AddMessage(c, s, "error", err.Error()); err != nil {
				c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
				return
			}
			c.Redirect(http.StatusFound, "/signin")
			return
		}

		redirectURI, err := usecases.NewAuthorization(cfg, db, s).Invoke(c, input.Email, input.Password)
		if err != nil {
			if usecaseErr, ok := err.(*cerrs.UsecaseError); ok {
				switch usecaseErr.Code {
				case http.StatusFound:
					if err := s.SetSessionData(c, "flushMessage", usecaseErr.Error()); err != nil {
						c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": usecaseErr.Error()})
						return
					}
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
