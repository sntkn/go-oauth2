package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal"
	"github.com/sntkn/go-oauth2/oauth2/internal/flashmessage"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

type CreateUserUsecase interface {
	Invoke(c *gin.Context, user *repository.User) error
}

type SignupInput struct {
	Name     string `form:"name" binding:"required"`
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func CreateUserHandler(c *gin.Context) {
	db, err := internal.GetFromContextIF[repository.OAuth2Repository](c, "db")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}
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

	uc := usecases.NewCreateUser(cfg, db, s)
	createUser(c, uc, s)
}

func createUser(c *gin.Context, uc CreateUserUsecase, s session.SessionClient) {
	var input SignupInput
	if err := c.ShouldBind(&input); err != nil {
		if flashErr := flashmessage.AddMessage(c, s, "error", err.Error()); flashErr != nil {
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": flashErr.Error()})
			return
		}
		c.Redirect(http.StatusFound, "/signup")
		return
	}

	user := &repository.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	if err := uc.Invoke(c, user); err != nil {
		if usecaseErr, ok := err.(*errors.UsecaseError); ok {
			switch usecaseErr.Code {
			case http.StatusBadRequest:
				if flashErr := flashmessage.AddMessage(c, s, flashmessage.Error, usecaseErr.Error()); flashErr != nil {
					c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": flashErr.Error()})
					return
				}
				c.Redirect(http.StatusFound, "/signup")
			case http.StatusInternalServerError:
				c.Error(errors.WithStack(err)) // TODO: trigger usecase
				c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": usecaseErr.Error()})
			}
			return
		}
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	if err := flashmessage.AddMessage(c, s, flashmessage.Success, "create user succeeded"); err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, "/signup-finished")
}
