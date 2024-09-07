package main

import (
	"log"

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

	dbConfig := &db.DBConfig{
		Host:     cfg.DBHost,
		Port:     uint16(cfg.DBPort),
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
	}

	database, err := db.Setup(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	injections := interfaces.NewInjection(database)

	routes.Setup(e, injections)

	// Start the server
	e.Logger.Fatal(e.Start(":18080"))
}

func responseMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		next := next(c)

		//if c.Response().Status == http.StatusInternalServerError {
		//	c.Response()
		//}
		return next
	}
}
