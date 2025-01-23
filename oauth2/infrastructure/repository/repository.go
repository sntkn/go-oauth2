package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
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
