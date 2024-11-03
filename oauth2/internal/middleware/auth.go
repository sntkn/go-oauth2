package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/accesstoken"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
)

func AuthMiddleware(cfg *config.Config, tokenParser accesstoken.Parser) gin.HandlerFunc {
	return func(c *gin.Context) {
		// "Authorization" ヘッダーを取得
		authHeader := c.GetHeader("Authorization")

		// "Authorization" ヘッダーが存在しない場合や、Bearer トークンでない場合はエラーを返す
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Missing or empty Authorization header")
			return
		}

		// "Bearer " のプレフィックスを取り除いてトークンを抽出
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := tokenParser.Parse(tokenStr, cfg.PublicKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set("claims", claims)
		c.Set("accessToken", tokenStr)

		c.Next()
	}
}
