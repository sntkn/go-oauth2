package user

import (
	"github.com/sntkn/go-oauth2/api/internal/modules/user/domain"
)

func NewUsecase(repo domain.Repository) *Usecase {
	return &Usecase{
		repository: repo,
	}
}

type Usecase struct {
	repository domain.Repository
}

func (s *Usecase) FindUser(id string) (*domain.User, error) {
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
