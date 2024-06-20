package user

import (
	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases"
)

func SignupFinishedHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		usecases.NewSignupFinished().Invoke(c)
	}
}
