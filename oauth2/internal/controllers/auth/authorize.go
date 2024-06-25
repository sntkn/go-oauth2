package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/pkg/redis"
	"github.com/sntkn/go-oauth2/oauth2/pkg/str"
)

type AuthorizeInput struct {
	ResponseType string `form:"response_type" binding:"required"`
	ClientID     string `form:"client_id" binding:"required"`
	Scope        string `form:"scope" binding:"required"`
	RedirectURI  string `form:"redirect_uri" binding:"required"`
	State        string `form:"state" binding:"required"`
}

func AuthrozeHandler(redisCli *redis.RedisCli, db *repository.Repository, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input AuthorizeInput

		if err := c.ShouldBind(&input); err != nil {
			c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
			return
		}

		if !str.IsValidUUID(input.ClientID) {
			err := fmt.Errorf("could not parse client_id")
			c.Error(errors.WithStack(err))
			c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
			return
		}

		s := session.NewSession(c, redisCli, cfg.SessionExpires)

		if err := usecases.NewAuthorize(redisCli, db, cfg, s).Invoke(c, input.ClientID, input.RedirectURI); err != nil {
			if usecaseErr, ok := err.(*cerrs.UsecaseError); ok {
				switch usecaseErr.Code {
				case http.StatusBadRequest:
					c.HTML(http.StatusFound, "400.html", gin.H{"error": err.Error()})
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
}
