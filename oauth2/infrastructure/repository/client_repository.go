package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

func NewClientRepository(db *sqlx.DB) IClientRepository {
	return ClientRepository{
		db: db,
	}
}

type Client struct {
	ID           uuid.UUID `db:"id"`
	Name         string    `db:"name"`
	RedirectURIs string    `db:"redirect_uris"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type IClientRepository interface {
	FindClientByClientID(clientID uuid.UUID) (*Client, error)
}

type ClientRepository struct {
	db *sqlx.DB
}

func (r ClientRepository) FindClientByClientID(clientID uuid.UUID) (*Client, error) {
	q := "SELECT id, redirect_uris FROM oauth2_clients WHERE id = $1"
	var c Client

	err := r.db.Get(&c, q, &clientID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &Client{}, errors.WithStack(err)
		}
		return nil, errors.WithStack(err)
	}
	return &c, errors.WithStack(err)
}
