package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sntkn/go-oauth2/oauth2/infrastructure/model"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

func NewDB(host, user, password, dbName string, port int) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	// PostgreSQLに接続
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return db, nil
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

type Repository struct {
	db *sqlx.DB
}

func (r *Repository) FindClientByClientID(clientID uuid.UUID) (*model.Client, error) {
	q := "SELECT id, redirect_uris FROM oauth2_clients WHERE id = $1"
	var c model.Client

	err := r.db.Get(&c, q, &clientID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &model.Client{}, errors.WithStack(err)
		}
		return nil, errors.WithStack(err)
	}
	return &c, errors.WithStack(err)
}

func (r Repository) FindUserByEmail(email string) (*model.User, error) {
	return nil, nil
}