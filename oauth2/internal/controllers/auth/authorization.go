package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/flashmessage"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/pkg/valkey"
)

//go:generate go run github.com/matryer/moq -out authorization_usecase_mock.go . AuthorizationUsecase
type AuthorizationUsecase interface {
	Invoke(input usecases.AuthorizationInput) (string, error)
}

type AuthorizationInput struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

type AuthorizationHandler struct {
	sessionManager session.SessionManager
	cfg            *config.Config
	uc             AuthorizationUsecase
}

func NewAuthorizationHandler(repo repository.OAuth2Repository, cfg *config.Config, valkeyCli valkey.ClientIF) *AuthorizationHandler {
	return &AuthorizationHandler{
		sessionManager: session.NewSessionManager(valkeyCli, cfg.SessionExpires),
		uc:             usecases.NewAuthorization(cfg, repo),
		cfg:            cfg,
	}
}

type SigninForm struct {
	Email string `form:"email"`
	Error string
}

func (h *AuthorizationHandler) Authorization(c *gin.Context) {
	sess := h.sessionManager.NewSession(c)
	var input AuthorizationInput

	if err := sess.SetNamedSessionData(c, "signin_form", SigninForm{
		Email: input.Email,
	}); err != nil {
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBind(&input); err != nil {
		if flashErr := flashmessage.AddMessage(c, sess, "error", err.Error()); flashErr != nil {
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": flashErr.Error()})
			return
		}
		c.Redirect(http.StatusFound, "/signin")
		return
	}

	var d AuthorizeInput
	if err := sess.GetNamedSessionData(c, "auth", &d); err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	redirectURI, err := h.uc.Invoke(usecases.AuthorizationInput{
		Email:       input.Email,
		Password:    input.Password,
		Scope:       d.Scope,
		RedirectURI: d.RedirectURI,
		ClientID:    d.ClientID,
		Expires:     h.cfg.AuthCodeExpires,
	})

	if err != nil {
		if usecaseErr, ok := err.(*errors.UsecaseError); ok {
			switch usecaseErr.Code {
			case http.StatusBadRequest:
				if flashErr := flashmessage.AddMessage(c, sess, "error", usecaseErr.Error()); flashErr != nil {
					c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": flashErr.Error()})
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

	if err := sess.DelSessionData(c, "auth"); err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, redirectURI)
}
