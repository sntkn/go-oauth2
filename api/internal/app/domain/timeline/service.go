package timeline

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

func (s *Service) RecentlyTimeline(userIDs []UserID) ([]*Timeline, error) {
	return s.repository.RecentlyTimeline(userIDs)
}
