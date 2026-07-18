package posts

import (
	"regexp"
	"strings"

	"github.com/eriscoo/blog-backend/internal/application"
	"github.com/eriscoo/blog-backend/internal/domain"
)

var reSlug = regexp.MustCompile(`[^a-z0-9-]`)

type Service struct {
	repo application.PostRepository
}

func New(repo application.PostRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll() ([]domain.Post, error) {
	return s.repo.FindAll()
}

func (s *Service) GetAllPublished() ([]domain.Post, error) {
	return s.repo.FindAllPublished()
}

func (s *Service) GetAllPublishedPaginated(offset, limit int) ([]domain.Post, int, error) {
	return s.repo.FindAllPublishedPaginated(offset, limit)
}

func (s *Service) GetByCategory(categoryID int, offset, limit int) ([]domain.Post, int, error) {
	return s.repo.FindAllByCategory(categoryID, offset, limit)
}

func (s *Service) GetByTag(tagID int, offset, limit int) ([]domain.Post, int, error) {
	return s.repo.FindAllByTag(tagID, offset, limit)
}

func (s *Service) GetByID(id int) (*domain.Post, error) {
	return s.repo.FindByID(id)
}

func (s *Service) GetBySlug(slug string) (*domain.Post, error) {
	return s.repo.FindBySlug(slug)
}

func (s *Service) Create(post *domain.Post) error {
	if post.Slug == "" {
		post.Slug = slugify(post.Title)
	}
	return s.repo.Create(post)
}

func (s *Service) Update(post *domain.Post) error {
	return s.repo.Update(post)
}

func (s *Service) Delete(id int) error {
	return s.repo.Delete(id)
}

func slugify(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = reSlug.ReplaceAllString(slug, "")
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}
