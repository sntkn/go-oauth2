package signup

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UseCase struct{}

func NewUseCase() *UseCase {
	return &UseCase{}
}

func (u *UseCase) Run(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", nil)
}
