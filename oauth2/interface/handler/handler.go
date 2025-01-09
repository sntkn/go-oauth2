package handler

import (
	"github.com/gin-gonic/gin"
)

type SigninHandler struct {
}

func (h *SigninHandler) Signin(c *gin.Context) {}

// func (h *OAuthHandler) Authorize(ctx echo.Context) error {
// 	// 認可リクエストを処理
// }
//
// func (h *OAuthHandler) Token(ctx echo.Context) error {
// 	// トークンリクエストを処理
// }
