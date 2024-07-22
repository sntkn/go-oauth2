package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/internal/flashmessage"
)

func SignupFinishedHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		mess, err := flashmessage.GetMessage(c)
		if err != nil {
			c.Error(errors.WithStack(err)) // TODO: trigger usecase
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
			return
		}

		c.HTML(http.StatusOK, "signup_finished.html", gin.H{"mess": mess})
	}
}
