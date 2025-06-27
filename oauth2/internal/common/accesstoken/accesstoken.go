package accesstoken

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

//go:generate go run github.com/matryer/moq -out token_gerenrator_mock.go . Generator Parser
type Generator interface {
	Generate(p *TokenParams, privateKeyBase64 string) (string, error)
}

type Parser interface {
	Parse(tokenStr string, publicKeyBase64 string) (*CustomClaims, error)
}

func NewTokenService() *Token {
	return &Token{}
}

type Token struct{}

type TokenParams struct {
	UserID    uuid.UUID
	ClientID  uuid.UUID
	Scope     string
	ExpiresAt time.Time
}

func (t *Token) Generate(p *TokenParams, privateKeyBase64 string) (string, error) {
	// JWTのペイロード（クレーム）を設定
	claims := jwt.MapClaims{
		"user_id":   p.UserID.String(),
		"client_id": p.ClientID.String(),
		"scope":     p.Scope,
		"exp":       p.ExpiresAt.Unix(),
		"iat":       time.Now().Unix(),
	}

	// JWTトークンを作成
	token := jwt.NewWithClaims(&jwt.SigningMethodEd25519{}, claims)

	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		return "", errors.WithStack(err)
	}

	privateKey := ed25519.PrivateKey(privateKeyBytes)

	// プライベートキーを使ってトークンを署名
	accessToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return accessToken, nil
}

type CustomClaims struct {
	UserID    string `json:"user_id"`
	ClientID  string `json:"client_id"`
	Scope     string
	ExpiresAt time.Time
	jwt.StandardClaims
}

func (t *Token) Parse(tokenStr string, publicKeyBase64 string) (*CustomClaims, error) {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// 公開鍵を使ってJWTをパース
	parsedToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return ed25519.PublicKey(publicKeyBytes), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*CustomClaims)
	if !ok || !parsedToken.Valid {
		err := errors.New("Invalid token")
		return nil, errors.WithStack(err)
	}

	return claims, nil
}
