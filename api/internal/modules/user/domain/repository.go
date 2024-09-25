package domain

import "github.com/sntkn/go-oauth2/api/internal/infrastructure/db/model"

type Repository interface {
	FindByID(id string) (*model.User, error)
}
