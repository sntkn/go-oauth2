package user

type UserRepository interface {
	FindByID(id string) (*User, error)
}

func NewService(repo UserRepository) *Service {
	return &Service{
		repository: repo,
	}
}

type Service struct {
	repository UserRepository
}

func (s *Service) FindUser(id string) (*User, error) {
	return s.repository.FindByID(id)
}
