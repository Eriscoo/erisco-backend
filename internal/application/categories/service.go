package categories

import (
	"github.com/eriscoo/blog-backend/internal/application"
	"github.com/eriscoo/blog-backend/internal/domain"
)

type Service struct {
	repo application.CategoryRepository
}

func New(repo application.CategoryRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll() ([]domain.Category, error) {
	return s.repo.FindAll()
}

func (s *Service) Create(name string) (domain.Category, error) {
	return s.repo.Create(name)
}

func (s *Service) Update(id int, name string) error {
	return s.repo.Update(id, name)
}

func (s *Service) Delete(id int) error {
	return s.repo.Delete(id)
}
