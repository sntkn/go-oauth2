package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/internal"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

type AuthorizeUsecase interface {
	Invoke(c *gin.Context, clientID string, redirectURI string) error
}

type AuthorizeInput struct {
	ResponseType string `form:"response_type" binding:"required"`
	ClientID     string `form:"client_id" binding:"required,uuid"`
	Scope        string `form:"scope" binding:"required"`
	RedirectURI  string `form:"redirect_uri" binding:"required"`
	State        string `form:"state" binding:"required"`
}

func AuthorizeHandler(c *gin.Context) { //nolint:dupl
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

	authorizeUsecase := usecases.NewAuthorize(cfg, db)
	authorize(c, s, authorizeUsecase)
}

func authorize(c *gin.Context, s session.SessionClient, uc AuthorizeUsecase) {
	var input AuthorizeInput

	if err := c.ShouldBind(&input); err != nil {
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		return
	}

	if err := uc.Invoke(c, input.ClientID, input.RedirectURI); err != nil {
		if usecaseErr, ok := err.(*cerrs.UsecaseError); ok {
			switch usecaseErr.Code {
			case http.StatusBadRequest:
				c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
			case http.StatusInternalServerError:
				c.Error(errors.WithStack(err)) // TODO: trigger usecase
				c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": usecaseErr.Error()})
			}
		} else {
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		}
		return
	}

	// セッションデータを書き込む
	if err := s.SetNamedSessionData(c, "auth", &input); err != nil {
		c.Error(err)
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/signin")
}
