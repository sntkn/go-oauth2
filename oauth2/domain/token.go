package domain

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

type TokenParams struct {
	AccessToken string
	ClientID    uuid.UUID
	UserID      uuid.UUID
	Scope       string
	ExpiresAt   time.Time
}

func NewToken(p TokenParams) Token {
	return &token{
		AccessToken: p.AccessToken,
		ClientID:    p.ClientID,
		UserID:      p.UserID,
		Scope:       p.Scope,
		ExpiresAt:   p.ExpiresAt,
	}
}

type StoreNewTokenParams struct {
	ClientID         uuid.UUID
	UserID           uuid.UUID
	Scope            string
	PrivateKeyBase64 string
	AdditionalMin    int
	Repo             TokenRepository
}

func StoreNewToken(p StoreNewTokenParams) (Token, error) {
	atoken := NewToken(TokenParams{
		ClientID: p.ClientID,
		UserID:   p.UserID,
		Scope:    p.Scope,
	})

	if err := atoken.SetNewAccessToken(p.PrivateKeyBase64); err != nil {
		return nil, err
	}
	atoken.SetNewExpiry(p.AdditionalMin)

	if err := p.Repo.StoreToken(atoken); err != nil {
		return nil, err
	}

	return atoken, nil
}

type Token interface {
	IsNotFound() bool
	GetAccessToken() string
	GetClientID() uuid.UUID
	GetUserID() uuid.UUID
	GetScope() string
	GetExpiresAt() time.Time
	SetNewAccessToken(privateKeyBase64 string) error
	Expiry() int64
	SetNewExpiry(additionalMin int)
}

type TokenRepository interface {
	StoreToken(Token) error
	FindToken(accessToken string) (Token, error)
	RevokeToken(accessToken string) error
}

type token struct {
	AccessToken string
	ClientID    uuid.UUID
	UserID      uuid.UUID
	Scope       string
	ExpiresAt   time.Time
}

func (t *token) IsNotFound() bool {
	return t.AccessToken == ""
}

func (t *token) GetAccessToken() string {
	return t.AccessToken
}

func (t *token) GetClientID() uuid.UUID {
	return t.ClientID
}

func (t *token) GetUserID() uuid.UUID {
	return t.UserID
}

func (t *token) GetScope() string {
	return t.Scope
}

func (t *token) GetExpiresAt() time.Time {
	return t.ExpiresAt
}

func (t *token) SetNewAccessToken(privateKeyBase64 string) error {
	newtoken, err := t.Generate(privateKeyBase64)
	if err != nil {
		return err
	}

	t.AccessToken = newtoken

	return nil
}

func (t *token) SetExpiresAt(tim time.Time) {
	t.ExpiresAt = tim
}

func (t *token) Expiry() int64 {
	return t.ExpiresAt.Unix()
}

func (t *token) SetNewExpiry(additionalMin int) {
	t.ExpiresAt = time.Now().Add(time.Duration(additionalMin) * time.Minute)
}

func (t *token) Generate(privateKeyBase64 string) (string, error) {
	// JWTのペイロード（クレーム）を設定
	claims := jwt.MapClaims{
		"user_id":   t.UserID.String(),
		"client_id": t.ClientID.String(),
		"scope":     t.Scope,
		"exp":       t.ExpiresAt.Unix(),
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

func (t *token) Parse(tokenStr string, publicKeyBase64 string) (*CustomClaims, error) {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// 公開鍵を使ってJWTをパース
	parsedToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
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
