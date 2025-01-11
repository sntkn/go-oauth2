package authentication

import (
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/infrastructure/model"
)

type IAuthenticationRepository interface {
	FindUserByEmail(email string) (*model.User, error)
	FindClientByClientID(clientID uuid.UUID) (*model.Client, error)
}
