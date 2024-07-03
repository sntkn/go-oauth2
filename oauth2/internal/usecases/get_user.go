package usecases

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/internal/accesstoken"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/pkg/redis"
)

type GetUser struct {
	redisCli *redis.RedisCli
	db       *repository.Repository
}

func NewGetUser(redisCli *redis.RedisCli, db *repository.Repository) *GetUser {
	return &GetUser{
		redisCli: redisCli,
		db:       db,
	}
}

func (u *GetUser) Invoke(c *gin.Context) (repository.User, error) {
	// "Authorization" ヘッダーを取得
	authHeader := c.GetHeader("Authorization")

	var user repository.User

	// "Authorization" ヘッダーが存在しない場合や、Bearer トークンでない場合はエラーを返す
	if authHeader == "" {
		// err := fmt.Errorf("missing or empty authorization header")
		// c.Error(errors.WithStack(err)).SetType(gin.ErrorTypePublic).SetMeta(http.StatusUnauthorized)
		return user, cerrs.NewUsecaseError(http.StatusUnauthorized, "missing or empty authorization header")
	}

	// "Bearer " のプレフィックスを取り除いてトークンを抽出
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := accesstoken.Parse(tokenStr)

	if err != nil {
		return user, cerrs.NewUsecaseError(http.StatusUnauthorized, err.Error())
	}

	// TODO: find user
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return user, cerrs.NewUsecaseError(http.StatusUnauthorized, err.Error())
	}

	user, err = u.db.FindUser(userID)
	if err != nil {
		return user, cerrs.NewUsecaseError(http.StatusUnauthorized, err.Error())
	}

	return user, nil
}
