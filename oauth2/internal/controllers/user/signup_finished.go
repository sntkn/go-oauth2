package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal"
	"github.com/sntkn/go-oauth2/oauth2/internal/flashmessage"
)

func SignupFinishedHandler(c *gin.Context) {
	mess, err := internal.GetFromContext[flashmessage.Messages](c, "flashMessages")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "signup_finished.html", gin.H{"mess": mess})
}
