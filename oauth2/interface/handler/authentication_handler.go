package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/infrastructure/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/entity"
	"github.com/sntkn/go-oauth2/oauth2/internal/flashmessage"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/usecase"
)

func NewAuthenticationHandler(opt HandlerOption) *AuthenticationHandler {
	repo := repository.NewRepository(opt.DB)
	uc := usecase.NewAuthenticationUsecase(repo)
	return &AuthenticationHandler{
		uc:      uc,
		session: opt.Session,
		config:  opt.Config,
	}
}

type AuthenticationHandler struct {
	uc      usecase.IAuthenticationUsecase
	session session.SessionManager
	config  *config.Config
}

type EntrySign struct {
	ResponseType string `form:"response_type" binding:"required"`
	ClientID     string `form:"client_id" binding:"required,uuid"`
	Scope        string `form:"scope" binding:"required"`
	RedirectURI  string `form:"redirect_uri" binding:"required"`
	State        string `form:"state" binding:"required"`
}

func (h *AuthenticationHandler) Entry(c *gin.Context) {
	sess := h.session.NewSession(c)

	var sign EntrySign

	if err := c.ShouldBindQuery(&sign); err != nil {
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		return
	}

	clientID, err := uuid.Parse(sign.ClientID)
	if err != nil {
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		return
	}

	_, err = h.uc.AuthenticateClient(clientID, sign.RedirectURI)
	if err != nil {
		handleError(c, sess, err)
		return
	}

	// セッションデータを書き込む
	if err := sess.SetNamedSessionData(c, "sign", &sign); err != nil {
		c.Error(err)
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/client/signin")
}

func (h *AuthenticationHandler) Signin(c *gin.Context) {
	sess := h.session.NewSession(c)
	mess, err := flashmessage.Flash(c, sess)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}
	var sign EntrySign
	var form entity.SessionSigninForm

	if err := sess.GetNamedSessionData(c, "sign", &sign); err != nil {
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		return
	}

	if sign.ClientID == "" {
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": "invalid client_id"})
		return
	}

	if err := sess.FlushNamedSessionData(c, "signin_form", &form); err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "signin.html", gin.H{"f": form, "mess": mess})
}

type PostSigninInput struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

func (h *AuthenticationHandler) PostSignin(c *gin.Context) {
	sess := h.session.NewSession(c)
	var input PostSigninInput
	var sign EntrySign

	if err := c.ShouldBind(&input); err != nil {
		if flashErr := flashmessage.AddMessage(c, sess, "error", err.Error()); flashErr != nil {
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": flashErr.Error()})
			return
		}
		c.Redirect(http.StatusFound, "/signin")
		return
	}

	if err := sess.GetNamedSessionData(c, "sign", &sign); err != nil {
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		return
	}

	user, err := h.uc.AuthenticateUser(input.Email, input.Password)

	if err != nil {
		handleError(c, sess, err)
		// サインインフォームをセッションに保存
		if err := sess.SetNamedSessionData(c, "signin_form", SigninForm{
			Email: input.Email,
		}); err != nil {
			c.Error(errors.WithStack(err))
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
			return
		}
		return
	}

	// signセッションを削除
	if err := sess.DelSessionData(c, "sign"); err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	// サインインフォームセッションを削除
	if err := sess.DelSessionData(c, "signin_form"); err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	// ログイン状態をセッションに保存
	if err := sess.SetNamedSessionData(c, "login", AuthedUser{
		Email:       input.Email,
		UserID:      user.ID.String(),
		ClientID:    sign.ClientID,
		RedirectURI: sign.RedirectURI,
		Scope:       sign.Scope,
		Expires:     h.config.AuthCodeExpires,
	}); err != nil {
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/oauth2/consent")
}
