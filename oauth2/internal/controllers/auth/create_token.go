package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/redis"
)

func CreateTokenHandler(redisCli *redis.RedisCli, db *repository.Repository, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		usecases.NewCreateToken(redisCli, db, cfg).Invoke(c)
	}
}
