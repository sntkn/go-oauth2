package api

import (
	"net/http"

	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/sntkn/go-oauth2/api/internal/interfaces"
	"github.com/sntkn/go-oauth2/api/internal/interfaces/response"
	"github.com/sntkn/go-oauth2/api/internal/modules/user"
	"github.com/sntkn/go-oauth2/api/internal/modules/user/domain"
	"github.com/sntkn/go-oauth2/api/internal/modules/user/registry"
)

type GetUserParams struct {
	ID string `param:"id"`
}

type GetUserResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type UserUsecase interface {
	FindUser(id string) (*domain.User, error)
}
type UserHandler struct {
	i  *interfaces.Injections
	uc UserUsecase
}

func NewUserHandler(injections *interfaces.Injections) *UserHandler {
	repo := registry.NewRepository(injections.DB)
	uc := user.NewUsecase(repo)

	return &UserHandler{
		i:  injections,
		uc: uc,
	}
}

func (h *UserHandler) GetUser(c echo.Context) error {
	params := new(GetUserParams)
	if err := c.Bind(params); err != nil {
		return response.APIResponse(c, http.StatusBadRequest, errors.Wrap("Invalid parameters", 0))
	}

	user, err := h.uc.FindUser(params.ID)
	if err != nil {
		return response.APIResponse(c, http.StatusInternalServerError, errors.Wrap("Failed to retrieve users", 0))
	}

	var userResponse GetUserResponse

	if err := copier.Copy(&userResponse, &user); err != nil {
		return response.APIResponse(c, http.StatusBadRequest, errors.Wrap("Cant copy response parameters", 0))
	}

	return response.APIResponse(c, http.StatusOK, &userResponse)
}
