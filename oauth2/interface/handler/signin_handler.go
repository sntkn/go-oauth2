package handler

import (
	"github.com/gin-gonic/gin"
)

func NewSigninHandler(opt HandlerOption) *SigninHandler {
	return &SigninHandler{
		opt: opt,
	}
}

type SigninHandler struct {
	opt HandlerOption
}

func (h *SigninHandler) Signin(c *gin.Context) {}

// func (h *OAuthHandler) Authorize(ctx echo.Context) error {
// 	// 認可リクエストを処理
// }
//
// func (h *OAuthHandler) Token(ctx echo.Context) error {
// 	// トークンリクエストを処理
// }
