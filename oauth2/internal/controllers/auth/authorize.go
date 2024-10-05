package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/pkg/redis"
)

//go:generate go run github.com/matryer/moq -out authorize_usecase_mock.go . AuthorizeUsecase
type AuthorizeUsecase interface {
	Invoke(clientID string, redirectURI string) error
}

type AuthorizeInput struct {
	ResponseType string `form:"response_type" binding:"required"`
	ClientID     string `form:"client_id" binding:"required,uuid"`
	Scope        string `form:"scope" binding:"required"`
	RedirectURI  string `form:"redirect_uri" binding:"required"`
	State        string `form:"state" binding:"required"`
}

type AuthorizeHandler struct {
	sessionManager session.SessionManager
	uc             AuthorizeUsecase
}

func NewAuthorizeHandler(repo repository.OAuth2Repository, cfg *config.Config, redisCli redis.RedisClient) *AuthorizeHandler {
	return &AuthorizeHandler{
		sessionManager: session.NewSessionManager(redisCli, cfg.SessionExpires),
		uc:             usecases.NewAuthorize(repo),
	}
}

func (h *AuthorizeHandler) Authorize(c *gin.Context) { //nolint:dupl // No need for commonization.
	sess := h.sessionManager.NewSession(c)
	var input AuthorizeInput

	if err := c.ShouldBind(&input); err != nil {
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.Invoke(input.ClientID, input.RedirectURI); err != nil {
		if usecaseErr, ok := err.(*errors.UsecaseError); ok {
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
	if err := sess.SetNamedSessionData(c, "auth", &input); err != nil {
		c.Error(err)
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/signin")
}
