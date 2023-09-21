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

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
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

type User struct {
	ID       uuid.UUID `db:"id"`
	Email    string    `db:"email"`
	Password string    `db:"password"`
	// 他のユーザー属性をここに追加
}

type Client struct {
	ID           uuid.UUID `db:"id"`
	Name         string    `db:"name"`
	RedirectURIs string    `db:"redirect_uris"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type Code struct {
	Code        string    `db:"code"`
	ClientID    uuid.UUID `db:"client_id"`
	UserID      uuid.UUID `db:"user_id"`
	Scope       string    `db:"scope"`
	RedirectURI string    `db:"redirect_uri"`
	ExpiresAt   time.Time `db:"expired_at"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

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
	Code      string `json:"code"`
	GrantType string `json:"grant_type"`
}

type TokenOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expiry       string `json:"expiry"`
}

var redisClient *redis.Client

func init() {
	// Redisクライアントの初期化
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redisのアドレスとポート番号に合わせて変更してください
		Password: "",               // Redisにパスワードが設定されている場合は設定してください
		DB:       0,                // データベース番号
	})

	// ピングしてRedis接続を確認
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
}

func main() {
	// Ginルーターの初期化
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	// セッションミドルウェアのセットアップ
	r.Use(SessionMiddleware())

	// PostgreSQLへの接続文字列を作成
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// PostgreSQLに接続
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// GETリクエストを受け取るエンドポイントの定義
	r.GET("/authorize", func(c *gin.Context) {

		// /authorize?response_type=code&client_id=550e8400-e29b-41d4-a716-446655440000&scope=read&redirect_uri=http%3A%2F%2Flocalhost%3A8000%2Fcallback&state=ok
		var input AuthorizeInput
		// Query ParameterをAuthorizeInputにバインド
		if err := c.BindQuery(&input); err != nil {
			c.HTML(http.StatusBadRequest, "Could not bind JSON", gin.H{"error": err.Error()})
			return
		}

		if input.ResponseType == "" {
			c.HTML(http.StatusBadRequest, "Invalid response type", nil)
			return
		}
		if input.ResponseType != "code" {
			c.HTML(http.StatusBadRequest, "Invalid response type", nil)
		}

		if input.ClientID == "" {
			c.HTML(http.StatusBadRequest, "Invalid client_id", nil)
			return
		}
		if IsValidUUID(input.ClientID) == false {
			c.HTML(http.StatusBadRequest, "Invalid client_id. UUID must be a valid UUID", nil)
			return
		}

		if input.Scope == "" {
			c.HTML(http.StatusBadRequest, "Invalid scope", nil)
			return
		}

		if input.RedirectURI == "" {
			c.HTML(http.StatusBadRequest, "Invalid redirect_uri", nil)
			return
		}

		if input.State == "" {
			c.HTML(http.StatusBadRequest, "Invalid state", nil)
			return
		}
		log.Printf("%+v\n", input)

		// check client
		query := "SELECT id, redirect_uris FROM oauth2_clients WHERE id = $1"
		var client Client

		err = db.QueryRow(query, input.ClientID).Scan(&client.ID, &client.RedirectURIs)
		if err != nil {
			if err == sql.ErrNoRows {
				c.HTML(http.StatusBadRequest, fmt.Sprintf("Could not Find Client: %s", input.ClientID), gin.H{"error": err.Error()})
			} else {
				c.HTML(http.StatusInternalServerError, "Internal Server Error", gin.H{"error": err.Error()})
			}
			return
		}

		if client.RedirectURIs != input.RedirectURI {
			c.HTML(http.StatusBadRequest, "Redirect URI does not match", nil)
			return
		}

		// セッションデータを書き込む
		d, err := json.Marshal(input)
		if err != nil {
			c.HTML(http.StatusBadRequest, "Could not marshal JSON", err)
			return
		}
		err = SetSessionData(c, d)
		if err != nil {
			c.HTML(http.StatusBadRequest, "Failed to set session data", err)
			return
		}

		c.HTML(http.StatusOK, "index.html", gin.H{"input": input})
	})

	r.POST("/authorization", func(c *gin.Context) {

		var input AuthorizationInput
		// リクエストのJSONデータをAuthorizationInputにバインド
		if err := c.Bind(&input); err != nil {
			c.HTML(http.StatusBadRequest, "Could not bind JSON", gin.H{"error": err.Error()})
			return
		}

		if input.Email == "" {
			// TODO: redirect to autorize with parameters
			c.HTML(http.StatusBadRequest, "Invalid email", nil)
			return
		}

		if input.Password == "" {
			// TODO: redirect to autorize with parameters
			c.HTML(http.StatusBadRequest, "Invalid email", nil)
			return
		}

		// validate user credentials
		query := "SELECT id, email, password FROM users WHERE email = $1"
		var user User

		err = db.QueryRow(query, input.Email).Scan(&user.ID, &user.Email, &user.Password)
		if err != nil {
			if err == sql.ErrNoRows {
				// TODO: redirect to autorize with parameters
				c.HTML(http.StatusBadRequest, "Could not Find User", gin.H{"error": err.Error()})
			} else {
				c.HTML(http.StatusInternalServerError, "Internal Server Error", gin.H{"error": err.Error()})
			}
			return
		}

		// パスワードを比較して認証
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
		if err != nil {
			// TODO: redirect to autorize with parameters
			c.HTML(http.StatusBadRequest, "Invalid password", gin.H{"error": err.Error()})
			return
		}

		sessionData, err := GetSessionData(c)
		if err != nil {
			c.HTML(http.StatusBadRequest, "Could not unmarshal session data", err)
			return
		}

		var d AuthorizeInput
		err = json.Unmarshal(sessionData, &d)
		if err != nil {
			c.HTML(http.StatusBadRequest, "Could not unmarshal session data", err)
			return
		}

		// create code
		expired := time.Now().AddDate(0, 0, 10)
		randomString, err := generateRandomString(32)
		if err != nil {
			c.HTML(http.StatusBadRequest, "Could not generate code generate random string", err)
			return
		}
		q := `
			INSERT INTO oauth2_codes
				(code, client_id, user_id, scope, redirect_uri, expires_at, created_at, updated_at)
			VALUES
				($1, $2, $3, $4, $5, $6, $7, $8)
		`
		_, err = db.Exec(q, randomString, d.ClientID, user.ID, d.Scope, d.RedirectURI, expired, time.Now(), time.Now())
		if err != nil {
			c.HTML(http.StatusBadRequest, "Could not create code: %v\n", err)
			return
		}

		// TODO: clear session data

		c.Redirect(http.StatusFound, fmt.Sprintf("%s?code=%s", d.RedirectURI, randomString))
	})

	r.POST("/token", func(c *gin.Context) {
		var input TokenInput
		if err := c.BindJSON(&input); err != nil {
			c.HTML(http.StatusBadRequest, "Could not bind JSON", gin.H{"error": err.Error()})
			return
		}

		// grant_type = authorization_code
		if input.GrantType != "authorization_code" {
			c.HTML(http.StatusForbidden, fmt.Sprintf("Invalid grant type: %s", input.GrantType), nil)
		}
		// code has expired
		query := "SELECT expires_at FROM oauth2_codes WHERE code = $1"
		var code Code

		err = db.QueryRow(query, input.Code).Scan(&code.ExpiresAt)
		if err != nil {
			if err == sql.ErrNoRows {
				// TODO: redirect to autorize with parameters
				c.HTML(http.StatusForbidden, "Could not Find Code", gin.H{"error": err.Error()})
			} else {
				c.HTML(http.StatusInternalServerError, "Internal Server Error", gin.H{"error": err})
			}
			return
		}
		currentTime := time.Now()
		if currentTime.After(code.ExpiresAt) {
			c.HTML(http.StatusForbidden, "Authorization Code expired", nil)
			return
		}

		// revoke code
		// create token and refresh token
		output := TokenOutput{
			AccessToken:  "token",
			RefreshToken: "refresh_token",
			Expiry:       "exp",
		}
		c.JSON(http.StatusOK, output)
	})

	// サーバーをポート8080で起動
	r.Run(":8080")
}

// セッションミドルウェアの定義
func SessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// セッションIDをクッキーから取得
		sessionID, err := c.Cookie("sessionID")
		if err != nil {
			// セッションIDがない場合は新しいセッションIDを生成
			sessionID = GenerateSessionID()
			// クッキーにセッションIDをセット
			c.SetCookie("sessionID", sessionID, 3600, "/", "localhost", false, true)
		}

		// Redisからセッションデータを取得
		sessionData, err := redisClient.Get(c, sessionID).Result()
		if err != nil {
			// セッションデータが存在しない場合は空のデータをセット
			sessionData = ""
		}

		// セッションデータをコンテキストにセット
		c.Set("sessionData", sessionData)

		// 次のハンドラを実行
		c.Next()
	}
}

// セッションIDを生成する関数
func GenerateSessionID() string {
	return time.Now().Format("20060102150405")
}

// セッションデータを取得する関数
func GetSessionData(c *gin.Context) ([]byte, error) {
	sessionID, err := c.Cookie("sessionID")
	if err != nil {
		return nil, err
	}
	return redisClient.Get(c, sessionID).Bytes()
}

// セッションデータをRedisに書き込む関数
func SetSessionData(c *gin.Context, sessionData any) error {
	// セッションIDをクッキーから取得
	sessionID, err := c.Cookie("sessionID")
	if err != nil {
		return err
	}

	// Redisにセッションデータを書き込み
	return redisClient.Set(c, sessionID, sessionData, 0).Err()
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
