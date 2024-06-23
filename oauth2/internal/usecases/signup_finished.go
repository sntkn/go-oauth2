package usecases

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SignupFinished struct{}

func NewSignupFinished() *SignupFinished {
	return &SignupFinished{}
}

func (*SignupFinished) Invoke(c *gin.Context) {
	c.HTML(http.StatusOK, "signup_finished.html", nil)
}
