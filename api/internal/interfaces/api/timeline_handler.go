package api

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sntkn/go-oauth2/api/internal/domain/timeline"
	"github.com/sntkn/go-oauth2/api/internal/interfaces"
)

type GetRecentlyTimelineParams struct {
	ID uuid.UUID `query:"id"`
}

type GetRecentlyTimelineResponse struct {
	PostTime    time.Time                     `json:"post_time"`
	Content     string                        `json:"content"`
	PostUser    UserResponse                  `json:"post_users"`
	Inpressions uint32                        `json:"inpressions"`
	Likes       []uuid.UUID                   `json:"likes"`
	Reposts     []GetRecentlyTimelineResponse `json:"reposts"`
}

type UserResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func GetRecentlyTimeline(i *interfaces.Injections, opts ...*interfaces.Ops) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := new(GetRecentlyTimelineParams)
		if err := c.Bind(params); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid parameters",
			})
		}

		repo := timeline.NewRepository(i.DB)
		s := timeline.NewService(repo)

		tl, err := s.RecentlyTimeline(timeline.UserID(params.ID))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to retrieve users",
			})
		}
		response := make([]GetRecentlyTimelineResponse, len(tl))
		for i, l := range tl {
			response[i] = GetRecentlyTimelineResponse{
				PostTime:    l.PostTime,
				Content:     l.Content,
				PostUser:    UserResponse{ID: uuid.UUID(l.PostUser.ID), Name: l.PostUser.Name},
				Inpressions: l.Inpressions,
				Likes:       []uuid.UUID{},
				Reposts:     []GetRecentlyTimelineResponse{},
			}
		}

		return c.JSON(http.StatusOK, response)
	}
}
