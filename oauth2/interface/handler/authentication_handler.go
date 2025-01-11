package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/infrastructure/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/entity"
	"github.com/sntkn/go-oauth2/oauth2/internal/flashmessage"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/usecase"
)

func NewAuthenticationHandler(opt HandlerOption) *AuthenticationHandler {
	repo := repository.NewRepository(opt.DB)
	uc := usecase.NewAuthenticationUsecase(repo)
	return &AuthenticationHandler{
		uc:      uc,
		session: opt.Session,
	}
}

type AuthenticationHandler struct {
	uc      usecase.IAuthenticationUsecase
	session session.SessionManager
}

type EntrySignInput struct {
	ResponseType string `form:"response_type" binding:"required"`
	ClientID     string `form:"client_id" binding:"required,uuid"`
	Scope        string `form:"scope" binding:"required"`
	RedirectURI  string `form:"redirect_uri" binding:"required"`
	State        string `form:"state" binding:"required"`
}

func (h *AuthenticationHandler) Entry(c *gin.Context) {
	sess := h.session.NewSession(c)

	var input EntrySignInput

	if err := c.ShouldBind(&input); err != nil {
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		return
	}

	clientID, err := uuid.Parse(input.ClientID)
	if err != nil {
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		return
	}

	_, err = h.uc.AuthenticateClient(clientID, input.RedirectURI)
	if err != nil {
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

func (h *AuthenticationHandler) Signin(c *gin.Context) {
	sess := h.session.NewSession(c)
	mess, err := flashmessage.Flash(c, sess)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}
	var input EntrySignInput
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
