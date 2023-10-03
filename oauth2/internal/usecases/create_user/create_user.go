package create_user

import (
	"encoding/json"
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/internal/redis"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
)

type SignupInput struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Password string `form:"password"`
}

type RegistrationForm struct {
	Name  string `form:"name"`
	Email string `form:"email"`
	Error string
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
	s := session.NewSession(c, u.redisCli)
	var input SignupInput
	// Query ParameterをAuthorizeInputにバインド
	if err := c.Bind(&input); err != nil {
		c.Error(errors.WithStack(err))
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	u.SetSessionData(c, s, RegistrationForm{
		Name:  input.Name,
		Email: input.Email,
	})

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
		//err := fmt.Errorf("User %s already exists", input.Email)
		c.Redirect(http.StatusFound, "/signup")
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

func (u *UseCase) SetSessionData(c *gin.Context, s *session.Session, v any) error {
	d, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return s.SetSessionData(c, "create_user_form", d)
}
