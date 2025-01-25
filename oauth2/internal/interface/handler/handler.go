package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/sntkn/go-oauth2/oauth2/internal/common/flashmessage"
	"github.com/sntkn/go-oauth2/oauth2/internal/common/session"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

type SigninForm struct {
	Email string `form:"email"`
	Error string
}

type AuthedUser struct {
	Name        string
	Email       string
	UserID      string
	ClientID    string
	RedirectURI string
	Scope       string
	Expires     int
}

type HandlerOption struct {
	Session session.SessionManager
	DB      *sqlx.DB
	Config  *config.Config
}

func handleError(c *gin.Context, sess session.SessionClient, err error) {
	if usecaseErr, ok := err.(*errors.UsecaseError); ok {
		if flashErr := flashmessage.AddMessage(c, sess, "error", usecaseErr.Error()); flashErr != nil {
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": flashErr.Error()})
			return
		}
		switch usecaseErr.Code {
		case http.StatusFound:
			c.Redirect(http.StatusFound, usecaseErr.RedirectURI)
		case http.StatusBadRequest:
			c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		case http.StatusInternalServerError:
			c.Error(errors.WithStack(err)) // TODO: trigger usecase
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": usecaseErr.Error()})
		}
	} else {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
	}
}
