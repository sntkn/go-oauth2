package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

type TokenInput struct {
	Code         string `json:"code" binding:"required_without=RefreshToken"`
	RefreshToken string `json:"refresh_token" binding:"required_without=Code"`
	GrantType    string `json:"grant_type" binding:"required"`
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

		if input.GrantType != "authorization_code" && input.GrantType != "refresh_token" {
			err := fmt.Errorf("invalid grant type: %s", input.GrantType)
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

		if input.RefreshToken == "" {
			err := fmt.Errorf("invalid refresh token")
			c.Error(errors.WithStack(err))
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
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
