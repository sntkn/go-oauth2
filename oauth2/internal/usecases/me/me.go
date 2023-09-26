package me

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/internal/redis"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
)

var secretKey = []byte("test")

type TokenParams struct {
	UserID    uuid.UUID
	ClientID  uuid.UUID
	Scope     string
	ExpiresAt time.Time
}
type CustomClaims struct {
	UserID    string `json:"user_id"`
	ClientID  string `json:"client_id"`
	Scope     string
	ExpiresAt time.Time
	jwt.StandardClaims
}

type UseCase struct {
	redisCli *redis.RedisCli
	db       *repository.Repository
}

func NewUseCase(redisCli *redis.RedisCli, db *repository.Repository) *UseCase {
	return &UseCase{
		redisCli: redisCli,
		db:       db,
	}
}

func (u *UseCase) Run(c *gin.Context) {
	// "Authorization" ヘッダーを取得
	authHeader := c.GetHeader("Authorization")

	// "Authorization" ヘッダーが存在しない場合や、Bearer トークンでない場合はエラーを返す
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or empty Authorization header"})
		return
	}

	// "Bearer " のプレフィックスを取り除いてトークンを抽出
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	// JWTトークンをパース
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// シークレットキーまたは公開鍵を返すことが必要です
		return secretKey, nil
	})

	if err != nil {
		fmt.Println("JWTパースエラー:", err)
		return
	}

	// カスタムクレームを取得
	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		fmt.Println("JWTが無効です")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Token"})
		return
	}

	// TODO: find user
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid UserID"})
		return
	}

	user, err := u.db.FindUser(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	// TODO: response user
	c.JSON(http.StatusOK, user)
}
