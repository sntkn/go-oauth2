package main

import (
	"context"
	"fmt"
	"log"

	"github.com/sntkn/go-oauth2/oauth2/internal/redis"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases/authorization"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases/authorize"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases/delete_token"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases/me"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases/token"

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

	r.GET("/authorize", authorize.NewUseCase(redisCli, db).Run)
	r.POST("/authorization", authorization.NewUseCase(redisCli, db).Run)
	r.POST("/token", token.NewUseCase(redisCli, db).Run)
	r.GET("/me", me.NewUseCase(redisCli, db).Run)
	r.DELETE("/token", delete_token.NewUseCase(redisCli, db).Run)

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
