package routes

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/sntkn/go-oauth2/api/internal/interfaces"
	"github.com/sntkn/go-oauth2/api/internal/interfaces/api"
)

func Setup(e *echo.Echo, injections *interfaces.Injections) {
	h := api.NewHandler(injections)

	// Define the routes
	u := e.Group("/users")
	u.Use(echojwt.JWT([]byte("test")))
	u.GET("/:id", h.GetUser)
	tl := e.Group("/timeline")
	tl.Use(echojwt.JWT([]byte("test")))
	tl.GET("/timeline/:id", h.GetRecentlyTimeline)
}
