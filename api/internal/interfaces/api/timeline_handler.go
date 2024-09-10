package api

import (
	"net/http"
	"time"

	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/sntkn/go-oauth2/api/internal/domain/timeline"
	"github.com/sntkn/go-oauth2/api/internal/interfaces"
	"github.com/sntkn/go-oauth2/api/internal/interfaces/response"
)

type GetRecentlyTimelineParams struct {
	ID uuid.UUID `param:"id"`
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
			return response.APIResponse(c, http.StatusBadRequest, errors.Wrap("Invalid parameters", 0))
		}

		repo := timeline.NewRepository(i.DB)
		s := timeline.NewService(repo)

		tl, err := s.RecentlyTimeline(timeline.UserID(params.ID))
		if err != nil {
			return response.APIResponse(c, http.StatusBadRequest, errors.Wrap("Failed to retrieve timeline", 0))
		}

		var timelineResponse []GetRecentlyTimelineResponse

		if err := copier.Copy(&timelineResponse, &tl); err != nil {
			return response.APIResponse(c, http.StatusBadRequest, errors.Wrap("Cant copy response parameters", 0))
		}

		return response.APIResponse(c, http.StatusOK, &timelineResponse)
	}
}
