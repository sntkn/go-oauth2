package token

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/accesstoken"
	"github.com/sntkn/go-oauth2/oauth2/internal/logs"
	"github.com/sntkn/go-oauth2/oauth2/internal/redis"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
)

var secretKey = []byte("test")

type UseCase struct {
	redisCli *redis.RedisCli
	db       *repository.Repository
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

func NewUseCase(redisCli *redis.RedisCli, db *repository.Repository) *UseCase {
	return &UseCase{
		redisCli: redisCli,
		db:       db,
	}
}

func (u *UseCase) Run(c *gin.Context) {
	var input TokenInput
	if err := c.BindJSON(&input); err != nil {
		msg := "Error binding JSON input"
		logs.ErrorWithWrap(err, msg)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	// grant_type = authorization_code
	if input.GrantType != "authorization_code" && input.GrantType != "refresh_token" {
		msg := fmt.Sprintf("Invalid grant type: %s", input.GrantType)
		err := errors.New(msg)
		logs.ErrorWithStack(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	if input.GrantType == "authorization_code" {
		// code has expired
		code, err := u.db.FindValidOAuth2Code(input.Code, time.Now())
		if err != nil {
			if err == sql.ErrNoRows {
				// TODO: redirect to autorize with parameters
				msg := "Could not find oauth2 code"
				logs.ErrorWithWrap(err, msg)
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": msg})
			} else {
				msg := "FindValidOAuth2Code Excecution error"
				logs.ErrorWithWrap(err, msg)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": msg})
			}
			return
		}
		currentTime := time.Now()
		if currentTime.After(code.ExpiresAt) {
			msg := "Code has expired"
			err := errors.New(msg)
			logs.ErrorWithStack(err)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": msg})
			return
		}

		// create token and refresh token
		expiration := time.Now().Add(10 * time.Minute)
		t := accesstoken.TokenParams{
			UserID:    code.UserID,
			ClientID:  code.ClientID,
			Scope:     code.Scope,
			ExpiresAt: expiration,
		}
		token, err := accesstoken.Generate(t)
		if err != nil {
			msg := "Could not generate token"
			logs.ErrorWithWrap(err, msg)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		err = u.db.RegisterToken(repository.Token{
			AccessToken: token,
			ClientID:    code.ClientID,
			UserID:      code.UserID,
			Scope:       code.Scope,
			ExpiresAt:   expiration,
		})
		if err != nil {
			msg := "Could not register token"
			logs.ErrorWithWrap(err, msg)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		randomString, err := generateRandomString(32)
		refreshExpiration := time.Now().AddDate(0, 0, 10)
		if err != nil {
			msg := "Could not generate refresh token"
			logs.ErrorWithWrap(err, msg)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		err = u.db.RegesterRefreshToken(repository.RefreshToken{
			RefreshToken: randomString,
			AccessToken:  token,
			ExpiresAt:    refreshExpiration,
		})
		if err != nil {
			msg := "Could not register refresh token"
			logs.ErrorWithWrap(err, msg)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		// revoke code
		err = u.db.RevokeCode(input.Code)
		if err != nil {
			msg := "Could not revoke code"
			logs.ErrorWithWrap(err, msg)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		output := TokenOutput{
			AccessToken:  token,
			RefreshToken: randomString,
			Expiry:       expiration.Unix(),
		}
		c.JSON(http.StatusOK, output)
		return
	}

	// check paramters
	if input.RefreshToken == "" {
		msg := "Invalid refresh token"
		err := errors.New(msg)
		logs.ErrorWithStack(err)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": msg})
		return
	}
	// TODO: find refresh token, if not expired
	rt, err := u.db.FindValidRefreshToken(input.RefreshToken, time.Now())
	if err != nil {
		if err == sql.ErrNoRows {
			msg := "Could not find refresh token"
			logs.ErrorWithWrap(err, msg)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": msg})
		} else {
			msg := "FindValidRefreshToken execution error"
			logs.ErrorWithWrap(err, msg)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	// find access token

	// TODO: create token and refresh token
	tkn, err := u.db.FindToken(rt.AccessToken)
	if err != nil {
		if err == sql.ErrNoRows {
			// TODO: redirect to autorize with parameters
			msg := "Could not find token"
			logs.ErrorWithWrap(err, msg)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			msg := "FindToken execution error"
			logs.ErrorWithWrap(err, msg)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
		}
		return
	}
	expiration := time.Now().Add(10 * time.Minute)
	t := accesstoken.TokenParams{
		UserID:    tkn.UserID,
		ClientID:  tkn.ClientID,
		Scope:     tkn.Scope,
		ExpiresAt: expiration,
	}
	token, err := accesstoken.Generate(t)
	if err != nil {
		msg := "Could not generate token"
		logs.ErrorWithWrap(err, msg)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = u.db.RegisterToken(repository.Token{
		AccessToken: token,
		ClientID:    tkn.ClientID,
		UserID:      tkn.UserID,
		Scope:       tkn.Scope,
		ExpiresAt:   expiration,
	})
	if err != nil {
		msg := "Could not register token"
		logs.ErrorWithWrap(err, msg)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	randomString, err := generateRandomString(32)
	refreshExpiration := time.Now().AddDate(0, 0, 10)
	if err != nil {
		msg := "Could not generate refresh token"
		logs.ErrorWithWrap(err, msg)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = u.db.RegesterRefreshToken(repository.RefreshToken{
		RefreshToken: randomString,
		AccessToken:  token,
		ExpiresAt:    refreshExpiration,
	})
	if err != nil {
		msg := "Could not register refresh token"
		logs.ErrorWithWrap(err, msg)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// TODO: revoke old token and refresh token
	if err = u.db.RevokeToken(tkn.AccessToken); err != nil {
		msg := "Could not revoke token"
		logs.ErrorWithWrap(err, msg)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err = u.db.RevokeRefreshToken(rt.RefreshToken); err != nil {
		msg := "Could not revoke refresh token"
		logs.ErrorWithWrap(err, msg)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	output := TokenOutput{
		AccessToken:  token,
		RefreshToken: randomString,
		Expiry:       expiration.Unix(),
	}
	c.JSON(http.StatusOK, output)
}

func generateRandomString(length int) (string, error) {
	// ランダムなバイト列を生成
	randomBytes := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, randomBytes)
	if err != nil {
		return "", err
	}

	// URLセーフなBase64エンコード
	encodedString := base64.URLEncoding.EncodeToString(randomBytes)

	return encodedString, nil
}
