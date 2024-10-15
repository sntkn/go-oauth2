package user

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/flashmessage"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/internal/usecases"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/sntkn/go-oauth2/oauth2/pkg/valkey"
)

//go:generate go run github.com/matryer/moq -out create_user_usecase_mock.go . CreateUserUsecase
type CreateUserUsecase interface {
	Invoke(user *repository.User) error
}

type SignupInput struct {
	Name     string `form:"name" binding:"required"`
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func NewCreateUserHandler(repo repository.OAuth2Repository, cfg *config.Config, valkeyCli valkey.ClientIF) *CreateUserHandler {
	uc := usecases.NewCreateUser(repo)
	return &CreateUserHandler{
		sessionManager: session.NewSessionManager(valkeyCli, cfg.SessionExpires),
		uc:             uc,
	}
}

type RegistrationData struct {
	Name  string
	Email string
}

type CreateUserHandler struct {
	sessionManager session.SessionManager
	uc             CreateUserUsecase
}

func (h *CreateUserHandler) CreateUser(c *gin.Context) {
	sess := h.sessionManager.NewSession(c)

	var input SignupInput
	if err := c.ShouldBind(&input); err != nil {
		if flashErr := flashmessage.AddMessage(c, sess, "error", err.Error()); flashErr != nil {
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

	if err := sess.SetNamedSessionData(c, "signup_form", RegistrationData{
		Name:  user.Name,
		Email: user.Email,
	}); err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.Invoke(user); err != nil {
		if usecaseErr, ok := err.(*errors.UsecaseError); ok {
			switch usecaseErr.Code {
			case http.StatusBadRequest:
				if flashErr := flashmessage.AddMessage(c, sess, flashmessage.Error, usecaseErr.Error()); flashErr != nil {
					c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": flashErr.Error()})
					return
				}
				c.Redirect(http.StatusFound, "/signup")
				return
			default:
				c.Error(errors.WithStack(err)) // TODO: trigger usecase
				c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": usecaseErr.Error()})
				return
			}
		}
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	if err := sess.DelSessionData(c, "signup_form"); err != nil {
		log.Printf("Error deleting session data: %v", err)
	}

	if err := flashmessage.AddMessage(c, sess, flashmessage.Success, "create user succeeded"); err != nil {
		log.Printf("Error adding flash message: %v", err)
	}
	c.Redirect(http.StatusFound, "/signup-finished")
}
