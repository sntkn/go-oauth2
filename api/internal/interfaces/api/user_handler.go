package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sntkn/go-oauth2/api/internal/domain/user"
	"gorm.io/gorm"
)

func GetUser(c echo.Context) error {
	id := c.Param("id")
	db := c.Get("db").(*gorm.DB)

	s := user.NewService(db)

	user, err := s.FindUser(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve users",
		})
	}
	return c.JSON(http.StatusOK, user)
}
