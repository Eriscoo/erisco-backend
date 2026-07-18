package tags

import (
	"github.com/eriscoo/blog-backend/internal/application"
	"github.com/eriscoo/blog-backend/internal/domain"
)

type Service struct {
	repo application.TagRepository
}

func New(repo application.TagRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll() ([]domain.Tag, error) {
	return s.repo.FindAll()
}

func (s *Service) GetByName(name string) (*domain.Tag, error) {
	return s.repo.FindByName(name)
}

func (s *Service) Create(name string) (domain.Tag, error) {
	return s.repo.Create(name)
}

func (s *Service) Update(id int, name string) error {
	return s.repo.Update(id, name)
}

func (s *Service) Delete(id int) error {
	return s.repo.Delete(id)
}
