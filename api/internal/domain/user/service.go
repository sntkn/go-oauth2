package user

type Service struct {
	repository *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repository: repo,
	}
}

func (s *Service) FindUser(id string) (*User, error) {
	return s.repository.FindByID(id)
}
