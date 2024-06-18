package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/sntkn/go-oauth2/oauth2/internal/redis"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases/authorization"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases/authorize"
	createToken "github.com/sntkn/go-oauth2/oauth2/internal/usecases/create_token"
	createUser "github.com/sntkn/go-oauth2/oauth2/internal/usecases/create_user"
	deleteToken "github.com/sntkn/go-oauth2/oauth2/internal/usecases/delete_token"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases/me"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases/signin"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases/signup"
	signupFinished "github.com/sntkn/go-oauth2/oauth2/internal/usecases/signup_finished"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Ginルーターの初期化
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// エラーログを出力するミドルウェアを追加
	r.Use(ErrorLoggerMiddleware())

	config, err := GetEnv()
	if err != nil {
		slog.Error("Config Load Error", "message:", err)
		return
	}

	// Redis configuration
	redisCli, err := redis.NewClient(context.Background(), redis.Options{
		Addr:     "session:6379", // Redisのアドレスとポート番号に合わせて変更してください
		Password: "",             // Redisにパスワードが設定されている場合は設定してください
		DB:       0,              // データベース番号
	})
	if err != nil {
		slog.Error("Session Error", "message:", err)
		return
	}

	// PostgreSQLに接続
	db, err := repository.NewClient(repository.Conn{
		Host:     config.DBHost,
		Port:     uint32(config.DBPort),
		User:     config.DBUser,
		Password: config.DBPassword,
		DBName:   config.DBName,
	})
	if err != nil {
		slog.Error("Database Error", "message:", err)
		return
	}
	defer db.Close()

	r.GET("/signin", signin.NewUseCase(redisCli).Run)
	r.GET("/authorize", authorize.NewUseCase(redisCli, db).Run)
	r.POST("/authorization", authorization.NewUseCase(redisCli, db).Run)
	r.POST("/token", createToken.NewUseCase(redisCli, db).Run)
	r.GET("/me", me.NewUseCase(redisCli, db).Run)
	r.DELETE("/token", deleteToken.NewUseCase(redisCli, db).Run)
	r.GET("/signup", signup.NewUseCase(redisCli).Run)
	r.POST("/signup", createUser.NewUseCase(redisCli, db).Run)
	r.GET("/signup-finished", signupFinished.NewUseCase().Run)

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

		err := c.Errors.ByType(gin.ErrorTypePublic).Last()
		if err != nil {
			// 短縮して型アサーションとデフォルト値の設定を一行で
			statusCode := func() int {
				if sc, ok := err.Meta.(int); ok {
					return sc
				}
				return http.StatusInternalServerError
			}()
			c.JSON(statusCode, gin.H{"error": err.Error()})
		}
	}
}

func ping() (bool, error) {
	config, err := GetEnv()
	if err != nil {
		slog.Error("Config Load Error", "message:", err)
		return false, err
	}

	db, err := repository.NewClient(repository.Conn{
		Host:     config.DBHost,
		Port:     uint32(config.DBPort),
		User:     config.DBUser,
		Password: config.DBPassword,
		DBName:   config.DBName,
	})
	if err != nil {
		slog.Error("Database Error", "message:", err)
		return false, err
	}
	defer db.Close()

	return true, nil
}
