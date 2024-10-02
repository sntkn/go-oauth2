package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/entity"
	"github.com/sntkn/go-oauth2/oauth2/internal/flashmessage"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/redis"
)

type SigninUsecaser interface {
	Invoke(c *gin.Context) (entity.SessionSigninForm, error)
}

type OAuthHandler struct {
	db             repository.OAuth2Repository
	sessionManager session.SessionManager
}

func NewOAuthHandler(db repository.OAuth2Repository, cfg *config.Config, redisCli redis.RedisClient) *OAuthHandler {
	sess := session.NewSessionManager(redisCli, cfg.SessionExpires)
	return &OAuthHandler{
		db:             db,
		sessionManager: sess,
	}
}

func (h *OAuthHandler) SigninHandler(c *gin.Context) { //nolint:dupl // No need for commonization.
	sess := h.sessionManager.NewSession(c)
	mess, err := flashmessage.Flash(c, sess)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}
	var input AuthorizeInput
	var form entity.SessionSigninForm

	if err := sess.GetNamedSessionData(c, "auth", &input); err != nil {
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		return
	}

	if input.ClientID == "" {
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": "invalid client_id"})
		return
	}

	if err := sess.FlushNamedSessionData(c, "signin_form", &form); err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "signin.html", gin.H{"f": form, "mess": mess})
}
