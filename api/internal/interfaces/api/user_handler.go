package api

import (
	"net/http"

	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/sntkn/go-oauth2/api/internal/domain/user"
	"github.com/sntkn/go-oauth2/api/internal/interfaces/response"
)

type GetUserParams struct {
	ID string `param:"id"`
}

type GetUserResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

func (h *Handler) GetUser(c echo.Context) error {
	params := new(GetUserParams)
	if err := c.Bind(params); err != nil {
		return response.APIResponse(c, http.StatusBadRequest, errors.Wrap("Invalid parameters", 0))
	}

	repo := user.NewRepository(h.i.DB)

	s := user.NewService(repo)

	user, err := s.FindUser(params.ID)
	if err != nil {
		return response.APIResponse(c, http.StatusInternalServerError, errors.Wrap("Failed to retrieve users", 0))
	}

	var userResponse GetUserResponse

	if err := copier.Copy(&userResponse, &user); err != nil {
		return response.APIResponse(c, http.StatusBadRequest, errors.Wrap("Cant copy response parameters", 0))
	}

	return response.APIResponse(c, http.StatusOK, &userResponse)
}
