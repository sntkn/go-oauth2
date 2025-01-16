package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/infrastructure/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/usecase"
)

func NewAuthorizationHandler(opt HandlerOption) *AuthorizationHandler {
	repo := repository.NewRepository(opt.DB)
	uc := usecase.NewAuthorizationUsecase(repo)
	return &AuthorizationHandler{
		uc:      uc,
		session: opt.Session,
		config:  opt.Config,
	}
}

type AuthorizationHandler struct {
	uc      usecase.IAuthorizationUsecase
	session session.SessionManager
	config  *config.Config
}

func (h *AuthorizationHandler) Consent(c *gin.Context) {
	sess := h.session.NewSession(c)

	// ログインセッションを取得
	var authUser AuthedUser
	if err := sess.GetNamedSessionData(c, "login", &authUser); err != nil {
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	// ログインしていない場合は400エラー
	if authUser.ClientID == "" {
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": "client not found"})
		return
	}

	clientID, err := uuid.Parse(authUser.ClientID)
	if err != nil {
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		return
	}

	// クライアント情報を取得
	client, err := h.uc.Consent(clientID)
	if err != nil {
		handleError(c, sess, err)
		return
	}

	c.HTML(http.StatusOK, "consent.html", gin.H{"cli": client})
}

type ConcentForm struct {
	Agree bool `form:"agree" binding:"required"`
}

func (h *AuthorizationHandler) PostConsent(c *gin.Context) {
	sess := h.session.NewSession(c)

	var concentForm ConcentForm

	// ログインセッションを取得
	var authUser AuthedUser
	if err := sess.GetNamedSessionData(c, "login", &authUser); err != nil {
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	// ログインしていない場合は400エラー
	if authUser.ClientID == "" {
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": "client not found"})
		return
	}

	if err := c.ShouldBind(&concentForm); err != nil {
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		return
	}

	// 同意画面のビジネスロジックを書く
	code, err := h.uc.GenerateAuthorizationCode(usecase.GenerateAuthorizationCodeParams{
		UserID:      authUser.UserID,
		ClientID:    authUser.ClientID,
		Scope:       authUser.Scope,
		RedirectURI: authUser.RedirectURI,
		Expires:     authUser.Expires,
	})

	if err != nil {
		handleError(c, sess, err)
		return
	}

	// リダイレクト
	c.Redirect(http.StatusFound, code.GenerateRedirectURIWithCode())
}
