package handler

import (
	"github.com/sntkn/go-oauth2/oauth2/infrastructure/repository"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
	"github.com/sntkn/go-oauth2/oauth2/usecase"
)

func NewAuthorizationHandler(opt HandlerOption) *AuthorizationHandler {
	repo := repository.NewRepository(opt.DB)
	uc := usecase.NewAuthorizationUsecase(repo)
	return &AuthorizationHandler{
		uc:      uc,
		session: opt.Session,
	}
}

type AuthorizationHandler struct {
	uc      usecase.IAuthorizationUsecase
	session session.SessionManager
}
