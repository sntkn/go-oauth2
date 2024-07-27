package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/internal"
	"github.com/sntkn/go-oauth2/oauth2/internal/flashmessage"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

func SigninHandler(c *gin.Context) {
	s, err := internal.GetFromContext[session.SessionClient](c, "session")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	mess, err := internal.GetFromContext[flashmessage.Messages](c, "flashMessages")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	cfg, err := internal.GetFromContext[config.Config](c, "cfg")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	form, err := usecases.NewSignin(cfg, s).Invoke(c)
	if err != nil {
		if usecaseErr, ok := err.(*cerrs.UsecaseError); ok {
			switch usecaseErr.Code {
			case http.StatusBadRequest:
				c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
			case http.StatusInternalServerError:
				c.Error(errors.WithStack(err)) // TODO: trigger usecase
				c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": usecaseErr.Error()})
			}
			return
		}
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "signin.html", gin.H{"f": form, "mess": mess})
}
