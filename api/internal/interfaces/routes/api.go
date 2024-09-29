package routes

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/sntkn/go-oauth2/api/internal/interfaces"
	"github.com/sntkn/go-oauth2/api/internal/interfaces/api"
)

func Setup(e *echo.Echo, injections *interfaces.Injections) {
	// Define the routes
	p := e.Group("/")
	p.Use(echojwt.JWT([]byte("test")))
	p.GET("user/:id", api.NewUserHandler(injections).GetUser)
	p.GET("/timeline/:id", api.NewTimelineHandler(injections).GetRecentlyTimeline)
}
