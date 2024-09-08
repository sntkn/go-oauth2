package api

import (
	"net/http"

	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sntkn/go-oauth2/api/internal/domain/user"
	"github.com/sntkn/go-oauth2/api/internal/interfaces"
	"github.com/sntkn/go-oauth2/api/internal/interfaces/response"
)

type GetUserQueryParams struct {
	ID string `query:"id"`
}

type GetUserResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

func GetUser(i *interfaces.Injections, opts ...*interfaces.Ops) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := new(GetUserQueryParams)
		if err := c.Bind(params); err != nil {
			return response.APIResponse(c, http.StatusBadRequest, errors.Wrap("Invalid parameters", 0))
		}

		repo := user.NewRepository(i.DB)

		s := user.NewService(repo)

		user, err := s.FindUser(params.ID)
		if err != nil {
			return response.APIResponse(c, http.StatusInternalServerError, errors.Wrap("Failed to retrieve users", 0))
		}
		data := &GetUserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		}
		return response.APIResponse(c, http.StatusOK, data)
	}
}
