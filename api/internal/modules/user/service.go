package user

import (
	"github.com/sntkn/go-oauth2/api/internal/modules/user/domain"
)

func NewService(repo domain.Repository) *Service {
	return &Service{
		repository: repo,
	}
}

type Service struct {
	repository domain.Repository
}

func (s *Service) FindUser(id string) (*domain.User, error) {
	u, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	user, err := domain.NewUser(u)
	if err != nil {
		return nil, err
	}
	return user, nil
}
