package repository

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sntkn/go-oauth2/oauth2/domain"
	"github.com/sntkn/go-oauth2/oauth2/infrastructure/model"
)

func NewClientRepository(db *sqlx.DB) *ClientRepository {
	return &ClientRepository{
		db: db,
	}
}

type ClientRepository struct {
	db *sqlx.DB
}

func (r *ClientRepository) FindClientByClientID(ctx context.Context, clientID uuid.UUID) (domain.Client, error) {
	q := "SELECT id, redirect_uris FROM oauth2_clients WHERE id = $1"
	mapper := func(c model.Client) (domain.Client, error) {
		redirectURIs := []string{}
		if c.RedirectURIs != "" {
			redirectURIs = strings.Split(c.RedirectURIs, ",")
		}
		return domain.NewClient(domain.ClientParams{
			ID:           c.ID,
			RedirectURIs: redirectURIs,
		}), nil
	}

	client, ok, err := fetchAndMap[model.Client, domain.Client](ctx, r.db, q, mapper, clientID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return domain.NewClient(domain.ClientParams{}), nil
	}

	return client, nil
}
