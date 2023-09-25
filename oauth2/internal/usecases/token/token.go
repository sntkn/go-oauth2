package token

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

type TokenParams struct {
	UserID    uuid.UUID
	ClientID  uuid.UUID
	Scope     string
	ExpiresAt time.Time
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	// grant_type = authorization_code
	if input.GrantType != "authorization_code" && input.GrantType != "refresh_token" {
		err := fmt.Errorf("Invalid grant type: %s", input.GrantType)
		c.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err})
		return
	}

	if input.GrantType == "authorization_code" {
		// code has expired
		code, err := u.db.FindValidOAuth2Code(input.Code, time.Now())
		if err != nil {
			if err == sql.ErrNoRows {
				// TODO: redirect to autorize with parameters
				c.JSON(http.StatusForbidden, gin.H{"error": err})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			}
			return
		}
		currentTime := time.Now()
		if currentTime.After(code.ExpiresAt) {
			err := fmt.Errorf("Authorization Code expired")
			c.Error(err)
			c.JSON(http.StatusForbidden, gin.H{"error": err})
			return
		}

		// create token and refresh token
		expiration := time.Now().Add(10 * time.Minute)
		t := TokenParams{
			UserID:    code.UserID,
			ClientID:  code.ClientID,
			Scope:     code.Scope,
			ExpiresAt: expiration,
		}
		token, err := generateAccessToken(t)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
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
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		randomString, err := generateRandomString(32)
		refreshExpiration := time.Now().AddDate(0, 0, 10)
		if err != nil {
			c.JSON(http.StatusForbidden, err)
			return
		}
		err = u.db.RegesterRefreshToken(repository.RefreshToken{
			RefreshToken: randomString,
			AccessToken:  token,
			ExpiresAt:    refreshExpiration,
		})
		if err != nil {
			c.JSON(http.StatusForbidden, err)
			return
		}

		// revoke code
		err = u.db.RevokeCode(input.Code)
		if err != nil {
			c.JSON(http.StatusForbidden, err)
			return
		}

		output := TokenOutput{
			AccessToken:  token,
			RefreshToken: randomString,
			Expiry:       expiration.Unix(),
		}
		c.JSON(http.StatusOK, output)
	} else {
		// TODO: check paramters
		// TODO: find refresh token, if not expired
		// TODO: find access token
		// TODO: create token and refresh token
		// TODO: revoke old token and refresh token

		output := TokenOutput{
			AccessToken:  "token",
			RefreshToken: "refresh token",
			Expiry:       0,
		}
		c.JSON(http.StatusOK, output)
	}
}

func generateAccessToken(p TokenParams) (string, error) {
	// JWTのペイロード（クレーム）を設定
	claims := jwt.MapClaims{
		"user_id":   p.UserID.String(),
		"client_id": p.ClientID.String(),
		"scope":     p.Scope,
		"exp":       p.ExpiresAt.Unix(),
		"iat":       time.Now().Unix(),
	}

	// JWTトークンを作成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// シークレットキーを使ってトークンを署名
	accessToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return accessToken, nil
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
