package usecases

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

type DeleteToken struct {
	db *repository.Repository
}

func NewDeleteToken(db *repository.Repository) *DeleteToken {
	return &DeleteToken{
		db: db,
	}
}

func (u *DeleteToken) Invoke(c *gin.Context) error {
	// "Authorization" ヘッダーを取得
	authHeader := c.GetHeader("Authorization")

	// "Authorization" ヘッダーが存在しない場合や、Bearer トークンでない場合はエラーを返す
	if authHeader == "" {
		return cerrs.NewUsecaseError(http.StatusUnauthorized, "Missing or empty Authorization header")
	}

	// "Bearer " のプレフィックスを取り除いてトークンを抽出
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	if err := u.db.RevokeToken(tokenStr); err != nil {
		return cerrs.NewUsecaseError(http.StatusInternalServerError, err.Error())
		//c.Error(errors.WithStack(err)).SetType(gin.ErrorTypePublic).SetMeta(http.StatusInternalServerError)
	}

	return nil
}
