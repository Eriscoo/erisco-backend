package profile

import (
	"github.com/eriscoo/blog-backend/internal/application"
	"github.com/eriscoo/blog-backend/internal/domain"
)

type Service struct {
	repo application.UserProfileRepository
}

func New(repo application.UserProfileRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(userID int) (*domain.UserProfile, error) {
	return s.repo.FindByUserID(userID)
}

func (s *Service) Update(profile *domain.UserProfile) error {
	return s.repo.Upsert(profile)
}
