package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sntkn/go-oauth2/api/internal/domain/timeline"
	"gorm.io/gorm"
)

func GetRecentlyTimeline(c echo.Context) error {
	id := c.Param("id")
	db := c.Get("db").(*gorm.DB)

	s := timeline.NewService(db)

	userID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	user, err := s.RecentlyTimeline(timeline.UserID(userID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve users",
		})
	}
	return c.JSON(http.StatusOK, user)
}
