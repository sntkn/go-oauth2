package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/internal/common/accesstoken"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

//go:generate go run github.com/matryer/moq -out get_user_usecase_mock.go . GetUserUsecase
type GetUserUsecase interface {
	Invoke(userID uuid.UUID) (repository.User, error)
}

func NewGetUserHandler(repo repository.OAuth2Repository, cfg *config.Config) *GetUserHandler {
	uc := usecases.NewGetUser(repo)
	return &GetUserHandler{
		uc:  uc,
		cfg: cfg,
	}
}

type GetUserHandler struct {
	uc  GetUserUsecase
	cfg *config.Config
}

func (h *GetUserHandler) GetUser(c *gin.Context) {
	var user repository.User

	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "claims not found"})
		return
	}

	customClaims, ok := claims.(*accesstoken.CustomClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid claims type"})
		return
	}

	userID, err := uuid.Parse(customClaims.UserID)
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
