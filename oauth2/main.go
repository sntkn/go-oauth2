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
	"github.com/sntkn/go-oauth2/oauth2/internal/common/accesstoken"
	"github.com/sntkn/go-oauth2/oauth2/internal/controllers/auth"
	"github.com/sntkn/go-oauth2/oauth2/internal/controllers/user"
	"github.com/sntkn/go-oauth2/oauth2/internal/middleware"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/validation"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/pkg/valkey"

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

	valkeyCli, err := valkey.NewClient(context.Background(), valkey.Options{
		Addr: []string{"kvs:6379"},
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

	r.GET("/signin", auth.NewSigninHandler(cfg, valkeyCli).Signin)
	r.GET("/authorize", auth.NewAuthorizeHandler(db, cfg, valkeyCli).Authorize)
	r.POST("/authorization", auth.NewAuthorizationHandler(db, cfg, valkeyCli).Authorization)

	r.POST("/token", auth.NewCreateTokenHandler(db, cfg).CreateToken)

	r.GET("/signup", user.NewSignupHandler(cfg, valkeyCli).Signup)
	r.POST("/signup", user.NewCreateUserHandler(db, cfg, valkeyCli).CreateUser)
	r.GET("/signup-finished", user.NewSignupFinishedHandler(cfg, valkeyCli).SignupFinished)

	g := r.Group("", middleware.AuthMiddleware(cfg, accesstoken.NewTokenService()))
	g.GET("/me", user.NewGetUserHandler(db, cfg).GetUser)
	g.DELETE("/token", auth.NewDeleteTokenHandler(db).DeleteToken)

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
