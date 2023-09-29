package accesstoken

import (
	"time"

	"github.com/cockroachdb/errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
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

func Generate(p TokenParams) (string, error) {
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
		return "", errors.WithStack(err)
	}

	return accessToken, nil
}

func Parse(tokenStr string) (*CustomClaims, error) {
	// JWTトークンをパース
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// シークレットキーまたは公開鍵を返すことが必要です
		return secretKey, nil
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
