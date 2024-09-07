package timeline

type Service struct {
	repository *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repository: repo,
	}
}

func (s *Service) RecentlyTimeline(userID UserID) ([]*Timeline, error) {
	// TODO: get follow userID
	userIDs := []UserID{userID}
	return s.repository.RecentlyTimeline(userIDs)
}
