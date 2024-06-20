package user

import (
	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/redis"
)

func SignupHandler(redisCli *redis.RedisCli, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		usecases.NewSignup(redisCli, cfg).Invoke(c)
	}
}
