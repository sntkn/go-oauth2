package usecases

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/internal/accesstoken"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

type GetUser struct {
	db repository.OAuth2Repository
}

func NewGetUser(db repository.OAuth2Repository) *GetUser {
	return &GetUser{
		db: db,
	}
}

func (u *GetUser) Invoke(c *gin.Context) (repository.User, error) {
	// "Authorization" ヘッダーを取得
	authHeader := c.GetHeader("Authorization")

	var user repository.User

	// "Authorization" ヘッダーが存在しない場合や、Bearer トークンでない場合はエラーを返す
	if authHeader == "" {
		return user, errors.NewUsecaseError(http.StatusUnauthorized, "missing or empty authorization header")
	}

	// "Bearer " のプレフィックスを取り除いてトークンを抽出
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := accesstoken.Parse(tokenStr)

	if err != nil {
		return user, errors.NewUsecaseError(http.StatusUnauthorized, err.Error())
	}

	// TODO: find user
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return user, errors.NewUsecaseError(http.StatusUnauthorized, err.Error())
	}

	user, err = u.db.FindUser(userID)
	if err != nil {
		return user, errors.NewUsecaseError(http.StatusUnauthorized, err.Error())
	}

	return user, nil
}
