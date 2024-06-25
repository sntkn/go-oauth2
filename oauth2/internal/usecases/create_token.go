package usecases

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/accesstoken"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/redis"
	"github.com/sntkn/go-oauth2/oauth2/pkg/str"
)

type CreateToken struct {
	redisCli *redis.RedisCli
	db       *repository.Repository
	cfg      *config.Config
}

type TokenInput struct {
	Code         string `json:"code"`
	RefreshToken string `json:"refresh_token"`
	GrantType    string `json:"grant_type"`
}

type TokenOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expiry       int64  `json:"expiry"`
}

func NewCreateToken(redisCli *redis.RedisCli, db *repository.Repository, cfg *config.Config) *CreateToken {
	return &CreateToken{
		redisCli: redisCli,
		db:       db,
		cfg:      cfg,
	}
}

const (
	randomStringLen = 32
	day             = 24 * time.Hour
)

func (u *CreateToken) Invoke(c *gin.Context) {
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
		u.createTokenByCode(c, &input)
		return
	}

	// check parameters
	if input.RefreshToken == "" {
		err := fmt.Errorf("invalid refresh token")
		c.Error(errors.WithStack(err))
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	u.CreateTokenByRefereshToken(c, &input)
}

func (u *CreateToken) createTokenByCode(c *gin.Context, input *TokenInput) {
	code, err := u.db.FindValidOAuth2Code(input.Code, time.Now())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// TODO: redirect to autorize with parameters
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	currentTime := time.Now()
	if currentTime.After(code.ExpiresAt) {
		err = fmt.Errorf("code has expired")
		c.Error(errors.WithStack(err))
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// create token and refresh token
	expiration := time.Now().Add(u.cfg.AuthTokenExpiresMin * time.Minute)
	t := accesstoken.TokenParams{
		UserID:    code.UserID,
		ClientID:  code.ClientID,
		Scope:     code.Scope,
		ExpiresAt: expiration,
	}
	token, err := accesstoken.Generate(t)
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err = u.db.RegisterToken(&repository.Token{
		AccessToken: token,
		ClientID:    code.ClientID,
		UserID:      code.UserID,
		Scope:       code.Scope,
		ExpiresAt:   expiration,
	}); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	randomString, err := str.GenerateRandomString(randomStringLen)
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	refreshExpiration := time.Now().Add(u.cfg.AuthRefreshTokenExpiresDay * day)
	if err = u.db.RegesterRefreshToken(&repository.RefreshToken{
		RefreshToken: randomString,
		AccessToken:  token,
		ExpiresAt:    refreshExpiration,
	}); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// revoke code
	if err = u.db.RevokeCode(input.Code); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, TokenOutput{
		AccessToken:  token,
		RefreshToken: randomString,
		Expiry:       expiration.Unix(),
	})
}

func (u *CreateToken) CreateTokenByRefereshToken(c *gin.Context, input *TokenInput) {
	// TODO: find refresh token, if not expired
	rt, err := u.db.FindValidRefreshToken(input.RefreshToken, time.Now())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	// find access token

	// TODO: create token and refresh token
	tkn, err := u.db.FindToken(rt.AccessToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// TODO: redirect to autorize with parameters
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	expiration := time.Now().Add(u.cfg.AuthTokenExpiresMin * time.Minute)

	token, err := accesstoken.Generate(accesstoken.TokenParams{
		UserID:    tkn.UserID,
		ClientID:  tkn.ClientID,
		Scope:     tkn.Scope,
		ExpiresAt: expiration,
	})
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err = u.db.RegisterToken(&repository.Token{
		AccessToken: token,
		ClientID:    tkn.ClientID,
		UserID:      tkn.UserID,
		Scope:       tkn.Scope,
		ExpiresAt:   expiration,
	}); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	randomString, err := str.GenerateRandomString(randomStringLen)
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	refreshExpiration := time.Now().Add(u.cfg.AuthRefreshTokenExpiresDay * day)

	if err = u.db.RegesterRefreshToken(&repository.RefreshToken{
		RefreshToken: randomString,
		AccessToken:  token,
		ExpiresAt:    refreshExpiration,
	}); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// TODO: revoke old token and refresh token
	if err = u.db.RevokeToken(tkn.AccessToken); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err = u.db.RevokeRefreshToken(rt.RefreshToken); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, TokenOutput{
		AccessToken:  token,
		RefreshToken: randomString,
		Expiry:       expiration.Unix(),
	})
}
