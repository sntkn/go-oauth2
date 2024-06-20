package usecases

import (
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/pkg/redis"
)

type SignupInput struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Password string `form:"password"`
}

type CreateUser struct {
	redisCli *redis.RedisCli
	db       *repository.Repository
}

func NewCreateUser(redisCli *redis.RedisCli, db *repository.Repository) *CreateUser {
	return &CreateUser{
		redisCli: redisCli,
		db:       db,
	}
}

func (u *CreateUser) Invoke(c *gin.Context) {
	s := session.NewSession(c, u.redisCli)
	var input SignupInput
	// Query ParameterをAuthorizeInputにバインド
	if err := c.Bind(&input); err != nil {
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	if err := s.SetNamedSessionData(c, "signup_form", RegistrationForm{
		Name:  input.Name,
		Email: input.Email,
	}); err != nil {
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	if input.Name == "" {
		c.Redirect(http.StatusFound, "/signup")
		return
	}

	if input.Email == "" {
		c.Redirect(http.StatusFound, "/signup")
		return
	}

	if input.Password == "" {
		c.Redirect(http.StatusFound, "/signup")
		return
	}

	// check email is existing
	eu, err := u.db.ExistsUserByEmail(input.Email)
	if err != nil {
		c.Error(err)
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	} else if eu {
		c.Redirect(http.StatusFound, "/signup")
		return
	}

	user := &repository.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	if err := u.db.CreateUser(user); err != nil {
		c.Error(err)
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	if err := s.DelSessionData(c, "signup_form"); err != nil {
		c.Error(err)
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/signup-finished")
}
