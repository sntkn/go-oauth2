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

	"github.com/gin-gonic/gin"

	repo "github.com/sntkn/go-oauth2/oauth2/infrastructure/repository"
	"github.com/sntkn/go-oauth2/oauth2/interface/handler"
	"github.com/sntkn/go-oauth2/oauth2/interface/presenter/bindings"
	"github.com/sntkn/go-oauth2/oauth2/internal/accesstoken"
	"github.com/sntkn/go-oauth2/oauth2/internal/controllers/auth"
	"github.com/sntkn/go-oauth2/oauth2/internal/controllers/user"
	"github.com/sntkn/go-oauth2/oauth2/internal/middleware"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/pkg/valkey"
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

	sdb, err := repo.NewDB(cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
	if err != nil {
		logger.Error("Database Error", "message:", err)
		return
	}
	defer sdb.Close()

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	if err := bindings.Setup(); err != nil {
		logger.Error("register validation error", "message:", err)
		return
	}

	r.Use(ErrorLoggerMiddleware(logger))

	opt := handler.HandlerOption{
		DB:      sdb,
		Session: session.NewSessionManager(valkeyCli, cfg.SessionExpires),
		Config:  cfg,
	}

	ah := handler.NewAuthenticationHandler(opt)
	r.GET("/client/sign-entry", ah.Entry)
	r.GET("/client/signin", ah.Signin)
	r.POST("/client/signin", ah.PostSignin)

	arh := handler.NewAuthorizationHandler(opt)
	r.GET("/oauth2/consent", arh.Consent)
	r.POST("/oauth2/consent", arh.PostConsent)

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
