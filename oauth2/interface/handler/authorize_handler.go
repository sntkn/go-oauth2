package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/interface/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/usecase"
)

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
	// sess := h.session.NewSession(c)
}

// func (h *OAuthHandler) Authorize(ctx echo.Context) error {
// 	// 認可リクエストを処理
// }
//
// func (h *OAuthHandler) Token(ctx echo.Context) error {
// 	// トークンリクエストを処理
// }
