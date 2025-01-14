package handler

import (
	"github.com/jmoiron/sqlx"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/pkg/config"
)

type SigninForm struct {
	Email string `form:"email"`
	Error string
}

type AuthedUser struct {
	Name        string
	Email       string
	UserID      string
	ClientID    string
	RedirectURI string
	Scope       string
	Expires     int
}

type HandlerOption struct {
	Session session.SessionManager
	DB      *sqlx.DB
	Config  *config.Config
}
