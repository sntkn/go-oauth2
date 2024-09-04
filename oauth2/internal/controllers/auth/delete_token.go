package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

type DeleteTokenUsecaser interface {
	Invoke(c *gin.Context) error
}

func DeleteTokenHandler(c *gin.Context) {
	db, err := internal.GetFromContextIF[repository.OAuth2Repository](c, "db")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}
	uc := usecases.NewDeleteToken(db)
	deleteToken(c, uc)
}

func deleteToken(c *gin.Context, uc DeleteTokenUsecaser) {
	if err := uc.Invoke(c); err != nil {
		if usecaseErr, ok := err.(*errors.UsecaseError); ok {
			c.AbortWithStatusJSON(usecaseErr.Code, gin.H{"error": usecaseErr.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}
