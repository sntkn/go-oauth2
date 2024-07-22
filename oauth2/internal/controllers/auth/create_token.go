package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

type TokenInput struct {
	Code         string `json:"code" binding:"required_without=RefreshToken,required_with_field_value=GrantType authorization_code"`
	RefreshToken string `json:"refresh_token" binding:"required_without=Code,required_with_field_value=GrantType refresh_token"`
	GrantType    string `json:"grant_type" binding:"required,oneof=authorization_code refresh_token"`
}

type TokenOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expiry       int64  `json:"expiry"`
}

func CreateTokenHandler(db *repository.Repository, cfg *config.Config) gin.HandlerFunc {

	return func(c *gin.Context) {
		var input TokenInput

		if err := c.BindJSON(&input); err != nil {
			c.Error(errors.WithStack(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if input.GrantType == "authorization_code" {
			token, err := usecases.NewCreateTokenByCode(cfg, db).Invoke(c, input.Code)
			if err != nil {
				if usecaseErr, ok := err.(*cerrs.UsecaseError); ok {
					c.AbortWithStatusJSON(usecaseErr.Code, gin.H{"error": usecaseErr.Error()})
					return
				}
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, TokenOutput{
				AccessToken:  token.AccessToken,
				RefreshToken: token.RefreshToken,
				Expiry:       token.Expiry,
			})
			return
		}

		token, err := usecases.NewCreateTokenByRefreshToken(cfg, db).Invoke(c, input.RefreshToken)
		if err != nil {
			if usecaseErr, ok := err.(*cerrs.UsecaseError); ok {
				c.AbortWithStatusJSON(usecaseErr.Code, gin.H{"error": usecaseErr.Error()})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, TokenOutput{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			Expiry:       token.Expiry,
		})
	}
}
