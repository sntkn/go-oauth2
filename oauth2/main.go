package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/sntkn/go-oauth2/oauth2/internal/controllers/auth"
	"github.com/sntkn/go-oauth2/oauth2/internal/controllers/user"
	"github.com/sntkn/go-oauth2/oauth2/internal/flashmessage"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/internal/validation"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/pkg/redis"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	_ "github.com/lib/pq"
)

const (
	readHederTimeout      = 5 * time.Second
	shutdownTimeoutSecond = 5 * time.Second
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.GetEnv()
	if err != nil {
		logger.Error("Config Load Error", "message:", err)
		return
	}

	redisCli, err := redis.NewClient(context.Background(), redis.Options{
		Addr:     "session:6379", // Redisのアドレスとポート番号に合わせて変更してください
		Password: "",             // Redisにパスワードが設定されている場合は設定してください
		DB:       0,              // データベース番号
	})
	if err != nil {
		logger.Error("Session Error", "message:", err)
		return
	}

	db, err := repository.NewClient(repository.Conn{
		Host:     cfg.DBHost,
		Port:     uint16(cfg.DBPort),
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
	})
	if err != nil {
		logger.Error("Database Error", "message:", err)
		return
	}
	defer db.Close()

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("required_with_field_value", validation.RequiredWithFieldValue); err != nil {
			logger.Error("register validation error")
		}
	}

	r.Use(ErrorLoggerMiddleware(logger))
	r.Use(func(c *gin.Context) {
		sess := session.NewSession(c, redisCli, cfg.SessionExpires)
		messages, _ := flashmessage.Flash(c, sess)
		c.Set("session", sess)
		c.Set("flashMessages", messages)
		c.Set("cfg", *cfg)
		c.Set("db", db)
	})

	r.GET("/signin", auth.SigninHandler)
	r.GET("/authorize", auth.AuthorizeHandler)
	r.POST("/authorization", auth.AuthorizationHandler)
	r.POST("/token", auth.CreateTokenHandler)
	r.DELETE("/token", auth.DeleteTokenHandler)
	r.GET("/me", user.GetUserHandler)
	r.GET("/signup", user.SignupHandler)
	r.POST("/signup", user.CreateUserHandler)
	r.GET("/signup-finished", user.SignupFinishedHandler)

	// サーバーの設定
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadHeaderTimeout: readHederTimeout,
	}

	// サーバーを非同期で起動
	go func() {
		if lasErr := srv.ListenAndServe(); lasErr != nil && lasErr != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", lasErr)
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
	if shutdownErr := srv.Shutdown(ctx); shutdownErr != nil {
		log.Printf("Server forced to shutdown: %+v\n", shutdownErr)
	}

	// ctx.Done() をキャッチする。5秒間のタイムアウト。
	<-ctx.Done()
	log.Println("timeout of 5 seconds.")

	log.Println("Server exiting")
}

// ErrorLoggerMiddleware はエラーログを出力するためのミドルウェアです。
func ErrorLoggerMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // 次のミドルウェアまたはハンドラを呼び出します

		for _, err := range c.Errors {
			logger.Error("errors occured", errors.LogStackTrace(err.Err))
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
		Port:     uint16(cfg.DBPort),
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
