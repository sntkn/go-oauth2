package authentication

import (
	"github.com/google/uuid"
)

type IAuthenticationRepository interface {
	FindUserByEmail(email string) (*User, error)
	FindClientByClientID(clientID uuid.UUID) (*Client, error)
}
