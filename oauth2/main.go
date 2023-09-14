package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/gin-gonic/gin"
)

type AuthorizeInput struct {
	ResponseType string `json:"response_type"`
	ClientId     string `json:"client_id"`
	Scope        string `json:"Scope"`
	RedirectURI  string `json:"redirect_uri"`
	State        string `json:"State"`
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

	// GETリクエストを受け取るエンドポイントの定義
	r.GET("/authorize", func(c *gin.Context) {

		// セッションからデータを取得
		sessionData := GetSessionData(c)
		fmt.Printf("%v\n", sessionData)

		// /authorize?response_type=code&client_id=abcdefg&scope=read&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fcallback&state=ok
		input := AuthorizeInput{}

		input.ResponseType = c.Query("response_type")
		log.Printf("Response type: %s", c.Query("response_type"))
		if input.ResponseType == "" {
			c.HTML(http.StatusBadRequest, "Invalid response type", nil)
			return
		}

		input.ClientId = c.Query("client_id")
		if input.ClientId == "" {
			c.HTML(http.StatusBadRequest, "Invalid client_id", nil)
			return
		}

		input.Scope = c.Query("scope")
		if input.Scope == "" {
			c.HTML(http.StatusBadRequest, "Invalid scope", nil)
			return
		}

		input.RedirectURI = c.Query("redirect_uri")
		if input.RedirectURI == "" {
			c.HTML(http.StatusBadRequest, "Invalid redirect_uri", nil)
			return
		}

		input.State = c.Query("state")
		if input.State == "" {
			c.HTML(http.StatusBadRequest, "Invalid state", nil)
			return
		}
		log.Printf("%+v\n", input)

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
func GetSessionData(c *gin.Context) string {
	sessionData, _ := c.Get("sessionData")
	return sessionData.(string)
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
