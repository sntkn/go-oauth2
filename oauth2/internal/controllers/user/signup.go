package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/internal"
	"github.com/sntkn/go-oauth2/oauth2/internal/entity"
	"github.com/sntkn/go-oauth2/oauth2/internal/flashmessage"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

type SignupUsecase interface {
	Invoke(c *gin.Context) (entity.SessionRegistrationForm, error)
}

type RegistrationForm struct {
	Name  string `form:"name"`
	Email string `form:"email"`
	Error string
}

func SignupHandler(c *gin.Context) {
	s, err := internal.GetFromContextIF[session.SessionClient](c, "session")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}
	cfg, err := internal.GetFromContext[config.Config](c, "cfg")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}
	mess, err := internal.GetFromContext[flashmessage.Messages](c, "flashMessages")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	uc := usecases.NewSignup(cfg, s)
	signup(c, mess, uc)
}

func signup(c *gin.Context, mess *flashmessage.Messages, uc SignupUsecase) {
	form, err := uc.Invoke(c)
	if err != nil {
		if usecaseErr, ok := err.(*cerrs.UsecaseError); ok {
			c.Error(errors.WithStack(err)) // TODO: trigger usecase
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": usecaseErr.Error()})
			return
		}
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "signup.html", gin.H{"f": form, "mess": mess})
}
