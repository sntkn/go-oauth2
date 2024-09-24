package domain

type Repository interface {
	RecentlyTimeline(userIDs []UserID) ([]*Timeline, error)
}
