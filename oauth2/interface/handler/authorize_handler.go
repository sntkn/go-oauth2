package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/infrastructure/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/usecase"
)

type AuthorizeInput struct {
	ResponseType string `form:"response_type" binding:"required"`
	ClientID     string `form:"client_id" binding:"required,uuid"`
	Scope        string `form:"scope" binding:"required"`
	RedirectURI  string `form:"redirect_uri" binding:"required"`
	State        string `form:"state" binding:"required"`
}

func NewAuthorizeHandler(opt HandlerOption) *AuthorizeHandler {
	userRepo := repository.NewUserRepository(opt.DB)
	clientRepo := repository.NewClientRepository(opt.DB)
	uc := usecase.NewAuthUsecase(userRepo, clientRepo)

	return &AuthorizeHandler{
		uc:      uc,
		session: opt.Session,
	}
}

type AuthorizeHandler struct {
	uc      usecase.IAuthUsecase
	session session.SessionManager
}

func (h *AuthorizeHandler) Authorize(c *gin.Context) {
	sess := h.session.NewSession(c)

	var input AuthorizeInput

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

// func (h *OAuthHandler) Authorize(ctx echo.Context) error {
// 	// 認可リクエストを処理
// }
//
// func (h *OAuthHandler) Token(ctx echo.Context) error {
// 	// トークンリクエストを処理
// }
