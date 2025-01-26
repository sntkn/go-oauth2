package domain

import (
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

func NewUser(p UserParams) *user {
	return &user{
		ID:        p.ID,
		Name:      p.Name,
		Email:     p.Email,
		Password:  p.Password,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

type User interface {
	GetID() uuid.UUID
	IsNotFound() bool
	IsPasswordMatch(password string) bool
}

type UserRepository interface {
	FindUserByEmail(email string) (User, error)
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
