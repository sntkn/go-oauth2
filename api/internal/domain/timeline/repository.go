package timeline

import (
	"log"

	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/api/internal/infrastructure/db/query"
	"gorm.io/gorm"
)

type Repository struct {
	query *query.Query
	gorm  *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		query: query.Use(db),
		gorm:  db,
	}
}

func (r *Repository) RecentlyTimeline(userIDs []UserID) ([]*Timeline, error) {
	postQuery := r.query.Post
	userIDStrings := make([]string, len(userIDs))

	for i, id := range userIDs {
		userIDStrings[i] = uuid.UUID(id).String()
	}

	posts, err := postQuery.Where(postQuery.UserID.In(userIDStrings...)).Order(postQuery.CreatedAt.Desc()).Find()
	if err != nil {
		log.Printf("Could not query posts: %v", err)
		return nil, err
	}

	var timeline []*Timeline

	for _, l := range posts {
		timeline = append(timeline, &Timeline{
			PostTime:    l.CreatedAt,
			Content:     l.Content,
			PostUser:    User{},
			Inpressions: 0,
			Likes:       []UserID{},
			Reposts:     []Timeline{},
		})
	}

	return timeline, nil
}
