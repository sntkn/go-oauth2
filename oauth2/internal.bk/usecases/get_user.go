package usecases

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

type GetUser struct {
	db repository.OAuth2Repository
}

func NewGetUser(db repository.OAuth2Repository) *GetUser {
	return &GetUser{
		db: db,
	}
}

func (u *GetUser) Invoke(userID uuid.UUID) (repository.User, error) {
	user, err := u.db.FindUser(userID)
	if err != nil {
		return user, errors.NewUsecaseError(http.StatusUnauthorized, err.Error())
	}

	return user, nil
}
