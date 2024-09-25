package timeline

import "github.com/sntkn/go-oauth2/api/internal/modules/timeline/domain"

type Usecase struct {
	repository domain.Repository
}

func NewUsecase(repo domain.Repository) *Usecase {
	return &Usecase{
		repository: repo,
	}
}

func (s *Usecase) RecentlyTimeline(userID domain.UserID) ([]*domain.Timeline, error) {
	// TODO: get follow userID
	userIDs := []domain.UserID{userID}
	posts, err := s.repository.RecentlyTimeline(userIDs)
	if err != nil {
		return nil, err
	}
	tl, err := domain.NewTimeline(posts)
	if err != nil {
		return nil, err
	}
	return tl, err
}
