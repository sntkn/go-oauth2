package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/internal"
	"github.com/sntkn/go-oauth2/oauth2/internal/flashmessage"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	cerrs "github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

type SignupInput struct {
	Name     string `form:"name" binding:"required"`
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func CreateUserHandler(c *gin.Context) {
	db, err := internal.GetFromContext[repository.SQLXOAuth2Repository](c, "db")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}
	s, err := internal.GetFromContext[session.Session](c, "session")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}
	cfg, err := internal.GetFromContext[config.Config](c, "cfg")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	var input SignupInput
	// Query ParameterをAuthorizeInputにバインド
	if err := c.Bind(&input); err != nil {
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		return
	}

	user := repository.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	if err := usecases.NewCreateUser(cfg, db, s).Invoke(c, user); err != nil {
		if usecaseErr, ok := err.(*cerrs.UsecaseError); ok {
			switch usecaseErr.Code {
			case http.StatusFound:
				if err := flashmessage.AddMessage(c, s, flashmessage.Error, usecaseErr.Error()); err != nil {
					c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
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
