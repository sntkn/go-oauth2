package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/redis"
)

func AuthrozationHandler(redisCli *redis.RedisCli, db *repository.Repository, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		usecases.NewAuthorization(redisCli, db, cfg).Invoke(c)
	}
}
