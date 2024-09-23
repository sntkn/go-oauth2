package timeline

type IRepository interface {
	RecentlyTimeline(userIDs []UserID) ([]*Timeline, error)
}

type Service struct {
	repository IRepository
}

func NewService(repo IRepository) *Service {
	return &Service{
		repository: repo,
	}
}

func (s *Service) RecentlyTimeline(userID UserID) ([]*Timeline, error) {
	// TODO: get follow userID
	userIDs := []UserID{userID}
	return s.repository.RecentlyTimeline(userIDs)
}
