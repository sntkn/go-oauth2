package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/entity"
	"github.com/sntkn/go-oauth2/oauth2/internal/flashmessage"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/redis"
)

type RegistrationForm struct {
	Name  string `form:"name"`
	Email string `form:"email"`
	Error string
}

func NewSignupHandler(cfg *config.Config, redisCli redis.RedisClient) *SignupHandler {
	return &SignupHandler{
		sessionManager: session.NewSessionManager(redisCli, cfg.SessionExpires),
	}
}

type SignupHandler struct {
	sessionManager session.SessionManager
}

func (h *SignupHandler) Signup(c *gin.Context) {
	sess := h.sessionManager.NewSession(c)
	mess, err := flashmessage.Flash(c, sess)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	var form entity.SessionRegistrationForm
	if err := sess.FlushNamedSessionData(c, "signup_form", &form); err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "signup.html", gin.H{"f": form, "mess": mess})
}
