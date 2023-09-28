package token

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/accesstoken"
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
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// grant_type = authorization_code
	if input.GrantType != "authorization_code" && input.GrantType != "refresh_token" {
		err := fmt.Errorf("Invalid grant type: %s", input.GrantType)
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.GrantType == "authorization_code" {
		// code has expired
		code, err := u.db.FindValidOAuth2Code(input.Code, time.Now())
		if err != nil {
			if err == sql.ErrNoRows {
				// TODO: redirect to autorize with parameters
				c.Error(err)
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			} else {
				c.Error(err)
				c.JSON(http.StatusInternalServerError, nil)
			}
			return
		}
		currentTime := time.Now()
		if currentTime.After(code.ExpiresAt) {
			err := fmt.Errorf("Authorization Code expired")
			c.Error(err)
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
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
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		randomString, err := generateRandomString(32)
		refreshExpiration := time.Now().AddDate(0, 0, 10)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		err = u.db.RegesterRefreshToken(repository.RefreshToken{
			RefreshToken: randomString,
			AccessToken:  token,
			ExpiresAt:    refreshExpiration,
		})
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		// revoke code
		err = u.db.RevokeCode(input.Code)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
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
		err := fmt.Errorf("Invalid refresh token")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: find refresh token, if not expired
	rt, err := u.db.FindValidRefreshToken(input.RefreshToken, time.Now())
	if err != nil {
		if err == sql.ErrNoRows {
			c.Error(err)
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	// find access token

	// TODO: create token and refresh token
	tkn, err := u.db.FindToken(rt.AccessToken)
	if err != nil {
		if err == sql.ErrNoRows {
			// TODO: redirect to autorize with parameters
			c.Error(err)
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	randomString, err := generateRandomString(32)
	refreshExpiration := time.Now().AddDate(0, 0, 10)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	err = u.db.RegesterRefreshToken(repository.RefreshToken{
		RefreshToken: randomString,
		AccessToken:  token,
		ExpiresAt:    refreshExpiration,
	})
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// TODO: revoke old token and refresh token
	if err = u.db.RevokeToken(tkn.AccessToken); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err = u.db.RevokeRefreshToken(rt.RefreshToken); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
