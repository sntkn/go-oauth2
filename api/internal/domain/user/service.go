package user

type Service struct {
	repository UserRepository
}

func NewService(repo UserRepository) *Service {
	return &Service{
		repository: repo,
	}
}

func (s *Service) FindUser(id string) (*User, error) {
	return s.repository.FindByID(id)
}
