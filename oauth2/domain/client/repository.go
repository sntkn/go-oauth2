package client

import (
	"github.com/google/uuid"
)

type ClientRepository interface {
	FindClientByClientID(clientID uuid.UUID) (*Client, error)
}
