package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

//go:generate go run github.com/matryer/moq -out delete_token_usecase_mock.go . DeleteTokenUsecase

type DeleteTokenUsecase interface {
	Invoke(tokenStr string) error
}

type DeleteTokenHandler struct {
	uc DeleteTokenUsecase
}

func NewDeleteTokenHandler(repo repository.OAuth2Repository) *DeleteTokenHandler {
	return &DeleteTokenHandler{
		uc: usecases.NewDeleteToken(repo),
	}
}

func (h *DeleteTokenHandler) DeleteToken(c *gin.Context) {
	token, exists := c.Get("accessToken")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "accessToken not found"})
		return
	}

	tokenStr, ok := token.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid accessToken type"})
		return
	}

	if err := h.uc.Invoke(tokenStr); err != nil {
		if usecaseErr, ok := err.(*errors.UsecaseError); ok {
			c.AbortWithStatusJSON(usecaseErr.Code, gin.H{"error": usecaseErr.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}
