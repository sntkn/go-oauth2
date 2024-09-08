package response

import (
	"net/http"

	"github.com/go-errors/errors"
	"github.com/labstack/echo/v4"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func APIResponse(c echo.Context, code int, message any) error {
	switch code {
	case http.StatusOK:
		response := Response{
			Status:  "success",
			Message: "",
			Data:    message,
		}
		return c.JSON(code, response)
	case http.StatusBadRequest:
		fallthrough
	case http.StatusForbidden:
		fallthrough
	case http.StatusInternalServerError:
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
	return c.JSON(http.StatusInternalServerError, "status not found")
}
