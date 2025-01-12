package handler

import (
	"github.com/jmoiron/sqlx"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
)

type SigninForm struct {
	Email string `form:"email"`
	Error string
}

type AuthedUser struct {
	Name        string
	Email       string
	ClientID    string
	RedirectURI string
}

type HandlerOption struct {
	Session session.SessionManager
	DB      *sqlx.DB
}
