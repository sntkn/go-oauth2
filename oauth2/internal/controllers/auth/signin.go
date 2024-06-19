package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/redis"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases"
)

func SigninHandler(redisCli *redis.RedisCli) gin.HandlerFunc {
	return func(c *gin.Context) {
		usecases.NewSignin(redisCli).Invoke(c)
	}
}
