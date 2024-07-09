package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SignupFinishedHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup_finished.html", nil)
	}
}
