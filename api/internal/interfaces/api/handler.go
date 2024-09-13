package api

import "github.com/sntkn/go-oauth2/api/internal/interfaces"

type Handler struct {
	i *interfaces.Injections
}

func NewHandler(injections *interfaces.Injections) *Handler {
	return &Handler{
		i: injections,
	}
}
