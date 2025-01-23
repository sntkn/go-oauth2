package usecases

import (
	"net/http"

	"github.com/sntkn/go-oauth2/oauth2/internal/repository"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

type DeleteToken struct {
	db repository.OAuth2Repository
}

func NewDeleteToken(db repository.OAuth2Repository) *DeleteToken {
	return &DeleteToken{
		db: db,
	}
}

func (u *DeleteToken) Invoke(tokenStr string) error {

	if err := u.db.RevokeToken(tokenStr); err != nil {
		return errors.NewUsecaseError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
