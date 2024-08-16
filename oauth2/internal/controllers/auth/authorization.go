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

type AuthorizationUsecase interface {
	Invoke(c *gin.Context, email string, password string) (string, error)
}

type AuthorizationInput struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

type SigninForm struct {
	Email string `form:"email"`
	Error string
}

func AuthorizationHandler(c *gin.Context) {
	db, err := internal.GetFromContextIF[repository.OAuth2Repository](c, "db")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}
	s, err := internal.GetFromContextIF[session.SessionClient](c, "session")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}
	cfg, err := internal.GetFromContext[config.Config](c, "cfg")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	uc := usecases.NewAuthorization(cfg, db, s)
	authorization(c, uc, s)
}

func authorization(c *gin.Context, uc AuthorizationUsecase, s session.SessionClient) {
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

	redirectURI, err := uc.Invoke(c, input.Email, input.Password)
	if err != nil {
		if usecaseErr, ok := err.(*cerrs.UsecaseError); ok {
			switch usecaseErr.Code {
			case http.StatusBadRequest:
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
