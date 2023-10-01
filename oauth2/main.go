package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/sntkn/go-oauth2/oauth2/internal/redis"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases/authorization"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases/authorize"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases/create_token"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases/create_user"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases/delete_token"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases/me"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases/signup"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases/signup_finished"

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

	slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// エラーログを出力するミドルウェアを追加
	r.Use(ErrorLoggerMiddleware())

	// Redis configuration
	redisCli, err := redis.NewClient(context.Background(), redis.Options{
		Addr:     "localhost:6379", // Redisのアドレスとポート番号に合わせて変更してください
		Password: "",               // Redisにパスワードが設定されている場合は設定してください
		DB:       0,                // データベース番号
	})
	if err != nil {
		slog.Error("Error: %v\n", err)
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
	r.POST("/token", create_token.NewUseCase(redisCli, db).Run)
	r.GET("/me", me.NewUseCase(redisCli, db).Run)
	r.DELETE("/token", delete_token.NewUseCase(redisCli, db).Run)
	r.GET("/signup", signup.NewUseCase().Run)
	r.POST("/signup", create_user.NewUseCase(redisCli, db).Run)
	r.GET("/signup-finished", signup_finished.NewUseCase().Run)

	// サーバーをポート8080で起動
	r.Run(":8080")
}

// ErrorLoggerMiddleware はエラーログを出力するためのミドルウェアです。
func ErrorLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // 次のミドルウェアまたはハンドラを呼び出します

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				slog.Error(fmt.Sprintf("%+v\n", err.Err))
			}
		}
	}
}
