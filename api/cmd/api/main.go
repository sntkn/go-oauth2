package main

import (
	"log"
	"math"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sntkn/go-oauth2/api/config"
	"github.com/sntkn/go-oauth2/api/internal/infrastructure/db"
	"github.com/sntkn/go-oauth2/api/internal/interfaces"
	"github.com/sntkn/go-oauth2/api/internal/interfaces/routes"
)

func main() {
	cfg, err := config.GetEnv()
	if err != nil {
		log.Fatal("could not get env:", err)
	}

	// バリデーション：int -> uint16 の変換前に範囲チェック
	if cfg.DBPort < 0 || cfg.DBPort > math.MaxUint16 {
		log.Fatalf("invalid DBPort: %d (must be 0..%d)", cfg.DBPort, math.MaxUint16)
	}
	dbPort := uint16(cfg.DBPort)

	dbConfig := &db.DBConfig{
		Host:     cfg.DBHost,
		Port:     dbPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
	}

	database, err := db.Setup(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	e := echo.New()
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Gzip())
	e.Use(middleware.Secure())
	injections := interfaces.NewInjection(database)

	routes.Setup(e, injections)

	// Start the server
	e.Logger.Fatal(e.Start(":18080"))
}
