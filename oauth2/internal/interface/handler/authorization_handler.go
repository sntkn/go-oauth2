package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/domain/authorization"
	"github.com/sntkn/go-oauth2/oauth2/infrastructure/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/common/accesstoken"
	"github.com/sntkn/go-oauth2/oauth2/internal/common/session"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecase"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

func NewAuthorizationHandler(opt HandlerOption) *AuthorizationHandler {
	repo := repository.NewAuthorizationRepository(opt.DB)
	tokenGen := accesstoken.NewTokenService()
	uc := usecase.NewAuthorizationUsecase(repo, opt.Config, tokenGen)
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

type TokenRequest struct {
	Code         string `json:"code" binding:"required_without=RefreshToken,required_with_field_value=GrantType authorization_code"`
	RefreshToken string `json:"refresh_token" binding:"required_without=Code,required_with_field_value=GrantType refresh_token"`
	GrantType    string `json:"grant_type" binding:"required,oneof=authorization_code refresh_token"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expiry       int64  `json:"expiry"`
}

func (h *AuthorizationHandler) Token(c *gin.Context) {
	var input TokenRequest
	var token *authorization.Token
	var err error

	if err := c.BindJSON(&input); err != nil {
		c.Error(errors.WithStack(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch input.GrantType {
	case "authorization_code":
		token, err = h.uc.GenerateTokenByCode(input.Code)
	case "refresh_token":
		token, err = h.uc.GenerateTokenByRefreshToken(input.RefreshToken)
	default:
		// ここには到達しない
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid grant type")})
		return
	}

	if err != nil {
		if usecaseErr, ok := err.(*errors.UsecaseError); ok {
			c.AbortWithStatusJSON(usecaseErr.Code, gin.H{"error": usecaseErr.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken.RefreshToken,
		Expiry:       token.Expiry(),
	})
}
