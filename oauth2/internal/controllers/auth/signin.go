package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases"
	"github.com/sntkn/go-oauth2/oauth2/pkg/redis"
)

func SigninHandler(redisCli *redis.RedisCli) gin.HandlerFunc {
	return func(c *gin.Context) {
		usecases.NewSignin(redisCli).Invoke(c)
	}
}
