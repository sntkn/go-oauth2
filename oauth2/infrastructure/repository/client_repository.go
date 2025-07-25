package repository

import (
	"database/sql"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sntkn/go-oauth2/oauth2/domain"
	"github.com/sntkn/go-oauth2/oauth2/infrastructure/model"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

func NewClientRepository(db *sqlx.DB) *ClientRepository {
	return &ClientRepository{
		db: db,
	}
}

type ClientRepository struct {
	db *sqlx.DB
}

func (r *ClientRepository) FindClientByClientID(clientID uuid.UUID) (domain.Client, error) {
	q := "SELECT id, redirect_uris FROM oauth2_clients WHERE id = $1"
	var c model.Client

	err := r.db.Get(&c, q, &clientID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.NewClient(domain.ClientParams{}), nil
		}
		return nil, errors.WithStack(err)
	}

	redirectURIs := []string{}
	if c.RedirectURIs != "" {
		redirectURIs = strings.Split(c.RedirectURIs, ",")
	}

	return domain.NewClient(domain.ClientParams{
		ID:           c.ID,
		RedirectURIs: redirectURIs,
	}), nil
}
