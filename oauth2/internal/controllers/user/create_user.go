package user

import (
	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/redis"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases"
)

func CreateUserHandler(redisCli *redis.RedisCli, db *repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		usecases.NewCreateUser(redisCli, db).Invoke(c)
	}
}
