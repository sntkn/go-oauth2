package deletetoken

import (
	"net/http"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/redis"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
)

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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or empty Authorization header"})
		return
	}

	// "Bearer " のプレフィックスを取り除いてトークンを抽出
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	if err := u.db.RevokeToken(tokenStr); err != nil {
		c.Error(errors.WithStack(err)).SetType(gin.ErrorTypePublic).SetMeta(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, nil)
}
