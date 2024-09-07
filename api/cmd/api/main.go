package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sntkn/go-oauth2/api/internal/infrastructure/db"
	"github.com/sntkn/go-oauth2/api/internal/interfaces"
	"github.com/sntkn/go-oauth2/api/internal/interfaces/routes"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	database, err := db.Setup()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	injections := interfaces.NewInjection(database)

	routes.Setup(e, injections)

	// Start the server
	e.Logger.Fatal(e.Start(":18080"))
}
