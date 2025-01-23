package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/common/accesstoken"
	"github.com/sntkn/go-oauth2/oauth2/internal/entity"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

//go:generate go run github.com/matryer/moq -out create_token_by_code_usecase_mock.go . CreateTokenByCodeUsecase
type CreateTokenByCodeUsecase interface {
	Invoke(authCode string) (*entity.AuthTokens, error)
}

//go:generate go run github.com/matryer/moq -out create_token_by_refresh_token_usecase_mock.go . CreateTokenByRefreshTokenUsecase
type CreateTokenByRefreshTokenUsecase interface {
	Invoke(refreshToken string) (*entity.AuthTokens, error)
}

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

type CreateTokenHandler struct {
	tokenUC        CreateTokenByCodeUsecase
	refreshTokenUC CreateTokenByRefreshTokenUsecase
}

func NewCreateTokenHandler(repo repository.OAuth2Repository, cfg *config.Config) *CreateTokenHandler {
	return &CreateTokenHandler{
		tokenUC:        usecases.NewCreateTokenByCode(cfg, repo, accesstoken.NewTokenService()),
		refreshTokenUC: usecases.NewCreateTokenByRefreshToken(cfg, repo, accesstoken.NewTokenService()),
	}
}

func (h *CreateTokenHandler) CreateToken(c *gin.Context) {
	var input TokenInput

	if err := c.BindJSON(&input); err != nil {
		c.Error(errors.WithStack(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch input.GrantType {
	case "authorization_code":
		token, err := h.tokenUC.Invoke(input.Code)
		if err != nil {
			if usecaseErr, ok := err.(*errors.UsecaseError); ok {
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
	case "refresh_token":
		token, err := h.refreshTokenUC.Invoke(input.RefreshToken)
		if err != nil {
			if usecaseErr, ok := err.(*errors.UsecaseError); ok {
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
	default:
		// ここには到達しない
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid grant type")})
	}
}
