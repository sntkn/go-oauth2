package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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

func CreateUserHandler(sessionCreator session.Creator, db *repository.Repository, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input SignupInput
		// Query ParameterをAuthorizeInputにバインド
		if err := c.Bind(&input); err != nil {
			c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
			return
		}

		s := sessionCreator(c)

		user := repository.User{
			Name:     input.Name,
			Email:    input.Email,
			Password: input.Password,
		}

		if err := usecases.NewCreateUser(cfg, db, s).Invoke(c, user); err != nil {
			if usecaseErr, ok := err.(*cerrs.UsecaseError); ok {
				switch usecaseErr.Code {
				case http.StatusFound:
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

		c.Redirect(http.StatusFound, "/signup-finished")
	}
}
