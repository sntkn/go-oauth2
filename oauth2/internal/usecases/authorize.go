package usecases

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/redis"
)

type AuthorizeInput struct {
	ResponseType string `form:"response_type" binding:"required"`
	ClientID     string `form:"client_id" binding:"required"`
	Scope        string `form:"scope" binding:"required"`
	RedirectURI  string `form:"redirect_uri" binding:"required"`
	State        string `form:"state" binding:"required"`
}

type Authorize struct {
	redisCli *redis.RedisCli
	db       *repository.Repository
	cfg      *config.Config
}

func NewAuthorize(redisCli *redis.RedisCli, db *repository.Repository, cfg *config.Config) *Authorize {
	return &Authorize{
		redisCli: redisCli,
		db:       db,
		cfg:      cfg,
	}
}

func (u *Authorize) Invoke(c *gin.Context) {
	s := session.NewSession(c, u.redisCli, u.cfg.SessionExpires)
	var input AuthorizeInput

	if err := c.ShouldBind(&input); err != nil {
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		return
	}

	if input.ResponseType != "code" {
		err := fmt.Errorf("invalid response_type: code must be 'code'")
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
	}

	if !IsValidUUID(input.ClientID) {
		err := fmt.Errorf("could not parse client_id")
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		return
	}

	// check client
	client, err := u.db.FindClientByClientID(input.ClientID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Error(err)
			c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		} else {
			c.Error(err)
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		}
		return
	}

	if client.RedirectURIs != input.RedirectURI {
		err = fmt.Errorf("redirect uri does not match")
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		return
	}

	// セッションデータを書き込む
	if err = s.SetNamedSessionData(c, "auth", &input); err != nil {
		c.Error(err)
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/signin")
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
