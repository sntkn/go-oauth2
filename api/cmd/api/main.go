package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sntkn/go-oauth2/api/internal/app/interfaces/db/model"
	"github.com/sntkn/go-oauth2/api/internal/app/interfaces/db/query"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	db := initDB()
	// defer db.Close()
	e.Use(dbMiddleware(db))

	// Define the routes
	e.GET("/users/:id", getUser)

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

func getUser(c echo.Context) error {
	id := c.Param("id")
	db := c.Get("db").(*gorm.DB)
	r := NewUserRepository(db)
	user, err := r.FindByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve users",
		})
	}
	return c.JSON(http.StatusOK, user)
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

type userRepository struct {
	query *query.Query
	gorm  *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		query: query.Use(db),
		gorm:  db,
	}
}

func (r *userRepository) FindByID(id string) (*model.User, error) {
	userQuery := r.query.User
	user, err := userQuery.Where(userQuery.ID.Eq(id)).First()

	if err != nil {
		log.Printf("Could not query user: %v", err)

		return nil, err
	}

	return user, nil
}
