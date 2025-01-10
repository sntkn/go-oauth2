package handler

import (
	"github.com/jmoiron/sqlx"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
)

type HandlerOption struct {
	Session session.SessionManager
	DB      *sqlx.DB
}
