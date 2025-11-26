package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserParams struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(p UserParams) User {
	return &user{
		ID:        p.ID,
		Name:      p.Name,
		Email:     p.Email,
		Password:  p.Password,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

//go:generate go run github.com/matryer/moq -out user_mock.go . User
type User interface {
	GetID() uuid.UUID
	IsNotFound() bool
	IsPasswordMatch(password string) bool
}

//go:generate go run github.com/matryer/moq -out user_repository_mock.go . UserRepository
type UserRepository interface {
	FindUserByEmail(ctx context.Context, email string) (User, error)
}

type user struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *user) GetID() uuid.UUID {
	return u.ID
}

func (u *user) IsNotFound() bool {
	return u.ID == uuid.Nil
}

func (u *user) IsPasswordMatch(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}
