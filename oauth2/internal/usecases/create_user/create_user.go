package create_user

import (
	"fmt"
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/internal/redis"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
)

type SignupInput struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Password string `form:"password"`
}

type UseCase struct {
	redisCli *redis.RedisCli
	db       *repository.Repository
}

func NewUseCase(redisCli *redis.RedisCli, db *repository.Repository) *UseCase {
	return &UseCase{
		redisCli: redisCli,
		db:       db,
	}
}

func (u *UseCase) Run(c *gin.Context) {
	var input SignupInput
	// Query ParameterをAuthorizeInputにバインド
	if err := c.Bind(&input); err != nil {
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		return
	}

	if input.Name == "" {
		err := fmt.Errorf("Invalid name")
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		return
	}

	if input.Email == "" {
		err := fmt.Errorf("Invalid email")
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		return
	}

	if input.Password == "" {
		err := fmt.Errorf("Invalid password")
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		return
	}

	// check email is existing
	eu, err := u.db.ExistsUserByEmail(input.Email)
	if err != nil {
		c.Error(err)
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	} else if eu {
		err := fmt.Errorf("User %s already exists", input.Email)
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusBadRequest, "400.html", gin.H{"error": err.Error()})
		return
	}

	user := repository.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	if err := u.db.CreateUser(user); err != nil {
		c.Error(err)
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/signup-finished")
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
