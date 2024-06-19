package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sntkn/go-oauth2/oauth2/internal/controllers/auth"
	"github.com/sntkn/go-oauth2/oauth2/internal/controllers/user"
	"github.com/sntkn/go-oauth2/oauth2/internal/redis"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	readHederTimeout      = 5 * time.Second
	shutdownTimeoutSecond = 5 * time.Second
)

func main() {
	// Ginルーターの初期化
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// エラーログを出力するミドルウェアを追加
	r.Use(ErrorLoggerMiddleware())

	cfg, err := config.GetEnv()
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
		Host:     cfg.DBHost,
		Port:     uint32(cfg.DBPort),
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
	})
	if err != nil {
		slog.Error("Database Error", "message:", err)
		return
	}
	defer db.Close()

	r.GET("/signin", auth.SigninHandler(redisCli))
	r.GET("/authorize", auth.AuthrozeHandler(redisCli, db))
	r.POST("/authorization", auth.AuthrozationHandler(redisCli, db, cfg))
	r.POST("/token", auth.CreateTokenHandler(redisCli, db))
	r.DELETE("/token", auth.DeleteTokenHandler(redisCli, db))
	r.GET("/me", user.GetUserHandler(redisCli, db))
	r.GET("/signup", user.SignupHandler(redisCli))
	r.POST("/signup", user.CreateUserHandler(redisCli, db))
	r.GET("/signup-finished", user.SignupFinishedHandler())

	// サーバーの設定
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadHeaderTimeout: readHederTimeout,
	}

	// サーバーを非同期で起動
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// シグナル受信のためのチャネルを作成
	quit := make(chan os.Signal, 1)
	// SIGINT（Ctrl+C）およびSIGTERMを受け取る
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// タイムアウト付きのコンテキストを設定
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeoutSecond)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %+v\n", err)
	}

	// ctx.Done() をキャッチする。5秒間のタイムアウト。
	<-ctx.Done()
	log.Println("timeout of 5 seconds.")

	log.Println("Server exiting")
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
	cfg, err := config.GetEnv()
	if err != nil {
		slog.Error("Config Load Error", "message:", err)
		return false, err
	}

	db, err := repository.NewClient(repository.Conn{
		Host:     cfg.DBHost,
		Port:     uint32(cfg.DBPort),
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
	})
	if err != nil {
		slog.Error("Database Error", "message:", err)
		return false, err
	}
	defer db.Close()

	return true, nil
}
