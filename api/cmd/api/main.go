package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sntkn/go-oauth2/api/internal/interfaces"
	"github.com/sntkn/go-oauth2/api/internal/interfaces/api"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	db := initDB()

	injections := interfaces.NewInjection(db)

	// Define the routes
	e.GET("/users/:id", api.GetUser(injections))
	e.GET("/timeline/:id", api.GetRecentlyTimeline(injections))

	// Start the server
	e.Logger.Fatal(e.Start(":18080"))
}

func dbMiddleware(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Store db in context
			c.Set("db", db)
			return next(c)
		}
	}
}

// Initialize DB connection
func initDB() *gorm.DB {
	dsn := "host=localhost user=admin password=admin dbname=auth port=5432 sslmode=disable TimeZone=Asia/Tokyo"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	return db
}
