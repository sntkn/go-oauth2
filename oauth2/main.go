package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/internal/redis"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "app"
	password = "pass"
	dbname   = "auth"
)

type AuthorizeInput struct {
	ResponseType string `form:"response_type"`
	ClientID     string `form:"client_id"`
	Scope        string `form:"scope"`
	RedirectURI  string `form:"redirect_uri"`
	State        string `form:"state"`
}

type AuthorizationInput struct {
	Email    string `form:"email"`
	Password string `form:"password"`
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

var secretKey = []byte("test")

func main() {
	// Ginルーターの初期化
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	// エラーログを出力するミドルウェアを追加
	r.Use(ErrorLoggerMiddleware())

	// Redis configuration
	redisCli, err := redis.NewClient(context.Background(), redis.Options{
		Addr:     "localhost:6379", // Redisのアドレスとポート番号に合わせて変更してください
		Password: "",               // Redisにパスワードが設定されている場合は設定してください
		DB:       0,                // データベース番号
	})
	if err != nil {
		log.Printf("%v", err)
		return
	}

	// PostgreSQLに接続
	db, err := repository.NewClient(repository.Conn{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbname,
	})
	defer db.Close()

	// GETリクエストを受け取るエンドポイントの定義
	r.GET("/authorize", func(c *gin.Context) {
		s := session.NewSession(c, redisCli)
		// /authorize?response_type=code&client_id=550e8400-e29b-41d4-a716-446655440000&scope=read&redirect_uri=http%3A%2F%2Flocalhost%3A8000%2Fcallback&state=ok
		var input AuthorizeInput
		// Query ParameterをAuthorizeInputにバインド
		if err := c.BindQuery(&input); err != nil {
			c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err})
			return
		}

		if input.ResponseType == "" {
			err := fmt.Errorf("Invalid response_type")
			c.Error(err)
			c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err})
			return
		}
		if input.ResponseType != "code" {
			err := fmt.Errorf("Invalid response_type: code must be 'code'")
			c.Error(err)
			c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err})
		}

		if input.ClientID == "" {
			err := fmt.Errorf("Invalid client_id")
			c.Error(err)
			c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err})
			return
		}
		if IsValidUUID(input.ClientID) == false {
			err := fmt.Errorf("Could not parse client_id")
			c.Error(err)
			c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err})
			return
		}

		if input.Scope == "" {
			err := fmt.Errorf("Invalid scope")
			c.Error(err)
			c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err})
			return
		}

		if input.RedirectURI == "" {
			err := fmt.Errorf("Invalid redirect_uri")
			c.Error(err)
			c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err})
			return
		}

		if input.State == "" {
			err := fmt.Errorf("Invalid state")
			c.Error(err)
			c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err})
			return
		}

		// check client
		client, err := db.FindClientByClientID(input.ClientID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err})
			} else {
				c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err})
			}
			return
		}

		if client.RedirectURIs != input.RedirectURI {
			err := fmt.Errorf("Redirect URI does not match")
			c.Error(err)
			c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err})
			return
		}

		// セッションデータを書き込む
		d, err := json.Marshal(input)
		if err != nil {
			c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err})
			return
		}
		err = s.SetSessionData(c, d)
		if err != nil {
			c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err})
			return
		}

		c.HTML(http.StatusOK, "index.html", gin.H{"input": input})
	})

	r.POST("/authorization", func(c *gin.Context) {
		s := session.NewSession(c, redisCli)
		var input AuthorizationInput
		// リクエストのJSONデータをAuthorizationInputにバインド
		if err := c.Bind(&input); err != nil {
			err := fmt.Errorf("Could not bind JSON")
			c.Error(err)
			c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err})
			return
		}

		if input.Email == "" {
			// TODO: redirect to autorize with parameters
			err := fmt.Errorf("Invalid email address")
			c.Error(err)
			c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err})
			return
		}

		if input.Password == "" {
			// TODO: redirect to autorize with parameters
			err := fmt.Errorf("Invalid password")
			c.Error(err)
			c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err})
			return
		}

		// validate user credentials
		user, err := db.FindUserByEmail(input.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				// TODO: redirect to autorize with parameters
				c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err})
			} else {
				c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err})
			}
			return
		}

		// パスワードを比較して認証
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
		if err != nil {
			// TODO: redirect to autorize with parameters
			c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err})
			return
		}

		sessionData, err := s.GetSessionData(c)
		if err != nil {
			c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err})
			return
		}

		var d AuthorizeInput
		err = json.Unmarshal(sessionData, &d)
		if err != nil {
			c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err})
			return
		}

		// create code
		expired := time.Now().AddDate(0, 0, 10)
		randomString, err := generateRandomString(32)
		if err != nil {
			c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err})
			return
		}

		clientID, err := uuid.Parse(d.ClientID)
		if err != nil {
			c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err})
			return
		}

		err = db.RegisterOAuth2Code(repository.Code{
			Code:        randomString,
			ClientID:    clientID,
			UserID:      user.ID,
			Scope:       d.Scope,
			RedirectURI: d.RedirectURI,
			ExpiresAt:   expired,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
		if err != nil {
			c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err})
			return
		}

		// TODO: clear session data

		c.Redirect(http.StatusFound, fmt.Sprintf("%s?code=%s", d.RedirectURI, randomString))
	})

	r.POST("/token", func(c *gin.Context) {
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
			code, err := db.FindValidOAuth2Code(input.Code, time.Now())
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

			err = db.RegisterToken(repository.Token{
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
			err = db.RegesterRefreshToken(repository.RefreshToken{
				RefreshToken: randomString,
				AccessToken:  token,
				ExpiresAt:    refreshExpiration,
			})
			if err != nil {
				c.JSON(http.StatusForbidden, err)
				return
			}

			// revoke code
			err = db.RevokeCode(input.Code)
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
	})

	// token refresh
	r.POST("/refresh", func(c *gin.Context) {

	})

	// サーバーをポート8080で起動
	r.Run(":8080")
}

// ErrorLoggerMiddleware はエラーログを出力するためのミドルウェアです。
func ErrorLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // 次のミドルウェアまたはハンドラを呼び出します

		// エラーが発生した場合はログに記録します
		if len(c.Errors) > 0 {
			for _, err := range c.Errors.Errors() {
				fmt.Printf("Error: %v\n", err)
			}
		}
	}
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
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
