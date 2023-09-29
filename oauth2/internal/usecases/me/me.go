package me

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/internal/accesstoken"
	"github.com/sntkn/go-oauth2/oauth2/internal/redis"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
)

var secretKey = []byte("test")

type UseCase struct {
	redisCli *redis.RedisCli
	db       *repository.Repository
}

func NewUseCase(redisCli *redis.RedisCli, db *repository.Repository) *UseCase {
	return &UseCase{
		redisCli: redisCli,
		db:       db,
	}
}

func (u *UseCase) Run(c *gin.Context) {
	// "Authorization" ヘッダーを取得
	authHeader := c.GetHeader("Authorization")

	// "Authorization" ヘッダーが存在しない場合や、Bearer トークンでない場合はエラーを返す
	if authHeader == "" {
		err := fmt.Errorf("Missing or empty Authorization header")
		c.Error(errors.WithStack(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// "Bearer " のプレフィックスを取り除いてトークンを抽出
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := accesstoken.Parse(tokenStr)

	if err != nil {
		c.Error(err)
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	// TODO: find user
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		c.Error(errors.WithStack(err))
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid UserID"})
		return
	}

	user, err := u.db.FindUser(userID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	// TODO: response user
	c.JSON(http.StatusOK, user)
}
