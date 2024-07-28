package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/internal"
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

func AuthrozationHandler(c *gin.Context) {
	db, err := internal.GetFromContext[repository.SQLXOAuth2Repository](c, "db")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}
	s, err := internal.GetFromContext[session.Session](c, "session")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}
	cfg, err := internal.GetFromContext[config.Config](c, "cfg")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

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
				if err := flashmessage.AddMessage(c, s, "error", usecaseErr.Error()); err != nil {
					c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
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
