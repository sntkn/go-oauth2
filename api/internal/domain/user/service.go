package user

import (
	"log"

	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/api/internal/domain/user/domain"
)

func NewService(repo domain.Repository) *Service {
	return &Service{
		repository: repo,
	}
}

type Service struct {
	repository domain.Repository
}

func (s *Service) FindUser(id string) (*domain.User, error) {
	u, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	userID, err := uuid.Parse(u.ID)
	if err != nil {
		log.Printf("UUIDの解析に失敗しました: %v", err)
		return nil, err
	}

	return &domain.User{
		ID:        userID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		// ここに他のフィールドを追加
	}, nil
}
