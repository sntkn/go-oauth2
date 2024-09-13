package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/sntkn/go-oauth2/api/internal/interfaces"
	"github.com/sntkn/go-oauth2/api/internal/interfaces/api"
)

func Setup(e *echo.Echo, injections *interfaces.Injections) {
	h := api.NewHandler(injections)

	// Define the routes
	e.GET("/users/:id", h.GetUser)
	e.GET("/timeline/:id", h.GetRecentlyTimeline)
}
