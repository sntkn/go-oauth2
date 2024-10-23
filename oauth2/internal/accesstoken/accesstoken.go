package accesstoken

import (
	"crypto/ed25519"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

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

func Generate(p TokenParams, privateKeyBase64 string) (string, error) {
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

func Parse(tokenStr string, publicKey string) (*CustomClaims, error) {
	// JWTトークンをパース
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(*jwt.Token) (any, error) {
		// シークレットキーまたは公開鍵を返すことが必要です
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	// カスタムクレームを取得
	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		err := errors.New("Invalid token")
		return nil, errors.WithStack(err)
	}

	return claims, nil
}
