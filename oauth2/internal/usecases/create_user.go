package usecases

import (
	"net/http"

	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

type RegistrationData struct {
	Name  string
	Email string
}

type CreateUser struct {
	db repository.OAuth2Repository
}

func NewCreateUser(db repository.OAuth2Repository) *CreateUser {
	return &CreateUser{
		db: db,
	}
}

func (u *CreateUser) Invoke(user *repository.User) error {

	// check email is existing
	eu, err := u.db.ExistsUserByEmail(user.Email)
	if err != nil {
		return errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	} else if eu {
		return errors.NewUsecaseError(http.StatusBadRequest, "input email already exists")
	}

	if err := u.db.CreateUser(user); err != nil {
		return errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
