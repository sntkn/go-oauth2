package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

func GetUserHandler(c *gin.Context) {
	db, err := internal.GetFromContext[repository.Repository](c, "db")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	user, err := usecases.NewGetUser(db).Invoke(c)
	if err != nil {
		if usecaseErr, ok := err.(*cerrs.UsecaseError); ok {
			c.AbortWithStatusJSON(usecaseErr.Code, gin.H{"error": usecaseErr.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// TODO: response user
	c.JSON(http.StatusOK, user)
}
