package repository

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sntkn/go-oauth2/oauth2/domain/model"
)

func NewClientRepository(db *sqlx.DB) IClientRepository {
	return ClientRepository{
		DB: db,
	}
}

type IClientRepository interface {
	FindClientByClientID(clientID uuid.UUID) (*model.Client, error)
}

type ClientRepository struct {
	DB *sqlx.DB
}

func (r ClientRepository) FindClientByClientID(clientID uuid.UUID) (*model.Client, error) {
	return nil, nil
}
