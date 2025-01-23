package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/common/flashmessage"
	"github.com/sntkn/go-oauth2/oauth2/internal/common/session"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/valkey"
)

func NewSignupFinishedHandler(cfg *config.Config, valkeyCli valkey.ClientIF) *SignupFinishedHandler {
	return &SignupFinishedHandler{
		sessionManager: session.NewSessionManager(valkeyCli, cfg.SessionExpires),
	}
}

type SignupFinishedHandler struct {
	sessionManager session.SessionManager
}

func (h *SignupFinishedHandler) SignupFinished(c *gin.Context) {
	sess := h.sessionManager.NewSession(c)
	mess, err := flashmessage.Flash(c, sess)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "signup_finished.html", gin.H{"mess": mess})
}
