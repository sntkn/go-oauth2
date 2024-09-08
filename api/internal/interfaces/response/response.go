package response

import (
	"net/http"

	"github.com/go-errors/errors"
	"github.com/labstack/echo/v4"
)

type Response[T any] struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func APIResponse[T any](c echo.Context, code int, message T) error {
	switch code {
	case http.StatusOK:
		response := Response[T]{
			Status:  "success",
			Message: "",
			Data:    message,
		}
		return c.JSON(code, response)

	case http.StatusBadRequest, http.StatusForbidden, http.StatusInternalServerError:
		err, ok := any(message).(*errors.Error)
		if !ok {
			return c.JSON(http.StatusInternalServerError, "could not bind error")
		}
		response := Response[any]{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		}
		return c.JSON(code, response)
	}
	return c.JSON(http.StatusInternalServerError, "status not found")
}
