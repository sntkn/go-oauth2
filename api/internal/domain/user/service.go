package user

import "gorm.io/gorm"

type Service struct {
	repository *Repository
}

func NewService(db *gorm.DB) *Service {
	r := NewRepository(db)

	return &Service{
		repository: r,
	}
}

func (s *Service) FindUser(id string) (*User, error) {
	return s.repository.FindByID(id)
}
