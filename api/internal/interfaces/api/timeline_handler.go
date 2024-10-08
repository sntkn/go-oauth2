package api

import (
	"net/http"
	"time"

	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/sntkn/go-oauth2/api/internal/interfaces"
	"github.com/sntkn/go-oauth2/api/internal/interfaces/response"
	"github.com/sntkn/go-oauth2/api/internal/modules/timeline"
	"github.com/sntkn/go-oauth2/api/internal/modules/timeline/domain"
	"github.com/sntkn/go-oauth2/api/internal/modules/timeline/registry"
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

type Usecase interface {
	RecentlyTimeline(userID domain.UserID) ([]*domain.Timeline, error)
}

type TimelineHandler struct {
	i  *interfaces.Injections
	uc Usecase
}

func NewTimelineHandler(injections *interfaces.Injections) *TimelineHandler {
	repo := registry.NewRepository(injections.DB)
	return &TimelineHandler{
		i:  injections,
		uc: timeline.NewUsecase(repo),
	}
}

func (h *TimelineHandler) GetRecentlyTimeline(c echo.Context) error {
	params := new(GetRecentlyTimelineParams)
	if err := c.Bind(params); err != nil {
		return response.APIResponse(c, http.StatusBadRequest, errors.Wrap("Invalid parameters", 0))
	}

	tl, err := h.uc.RecentlyTimeline(domain.UserID(params.ID))
	if err != nil {
		return response.APIResponse(c, http.StatusBadRequest, errors.Wrap("Failed to retrieve timeline", 0))
	}

	res := make([]GetRecentlyTimelineResponse, len(tl))

	if err := copier.Copy(&res, &tl); err != nil {
		return response.APIResponse(c, http.StatusBadRequest, errors.Wrap("Cant copy response parameters", 0))
	}

	return response.APIResponse(c, http.StatusOK, &res)
}
