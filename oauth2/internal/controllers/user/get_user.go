package user

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/internal/accesstoken"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

//go:generate go run github.com/matryer/moq -out get_user_usecase_mock.go . GetUserUsecase
type GetUserUsecase interface {
	Invoke(userID uuid.UUID) (repository.User, error)
}

func NewGetUserHandler(repo repository.OAuth2Repository, cfg *config.Config, tokenParser accesstoken.Parser) *GetUserHandler {
	uc := usecases.NewGetUser(repo)
	return &GetUserHandler{
		uc:          uc,
		cfg:         cfg,
		tokenParser: tokenParser,
	}
}

type GetUserHandler struct {
	uc          GetUserUsecase
	cfg         *config.Config
	tokenParser accesstoken.Parser
}

func (h *GetUserHandler) GetUser(c *gin.Context) {
	var user repository.User

	// "Authorization" ヘッダーを取得
	authHeader := c.GetHeader("Authorization")

	// "Authorization" ヘッダーが存在しない場合や、Bearer トークンでない場合はエラーを返す
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "Missing or empty Authorization header")
		return
	}

	// "Bearer " のプレフィックスを取り除いてトークンを抽出
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := h.tokenParser.Parse(tokenStr, h.cfg.PublicKey)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return

	}

	user, err = h.uc.Invoke(userID)
	if err != nil {
		if usecaseErr, ok := err.(*errors.UsecaseError); ok {
			c.AbortWithStatusJSON(usecaseErr.Code, gin.H{"error": usecaseErr.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// TODO: response user
	c.JSON(http.StatusOK, user)
}
