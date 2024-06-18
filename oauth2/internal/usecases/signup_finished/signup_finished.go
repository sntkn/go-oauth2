package signupFinished

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UseCase struct{}

func NewUseCase() *UseCase {
	return &UseCase{}
}

func (*UseCase) Run(c *gin.Context) {
	c.HTML(http.StatusOK, "signup_finished.html", nil)
}
