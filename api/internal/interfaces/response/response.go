package response

import (
	"net/http"

	"github.com/go-errors/errors"
	"github.com/labstack/echo/v4"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func APIResponse(c echo.Context, code int, message any) error {
	switch code {
	case http.StatusOK:
		return sendSuccessResponse(c, code, message)
	case http.StatusBadRequest, http.StatusForbidden, http.StatusInternalServerError:
		return sendErrorResponse(c, code, message)
	default:
		return c.JSON(http.StatusInternalServerError, "status not found")
	}
}

func sendSuccessResponse(c echo.Context, code int, data any) error {
	response := Response{
		Status:  "success",
		Message: "",
		Data:    data,
	}
	return c.JSON(code, response)
}

func sendErrorResponse(c echo.Context, code int, message any) error {
	err, ok := message.(*errors.Error)
	if !ok {
		return c.JSON(http.StatusInternalServerError, "could not bind error")
	}
	response := Response{
		Status:  "error",
		Message: err.Error(),
		Data:    nil,
	}
	return c.JSON(code, response)
}
