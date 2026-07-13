package application

import "github.com/eriscoo/blog-backend/internal/domain"

type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	FindByID(id int) (*domain.User, error)
}

type TagRepository interface {
	FindAll() ([]domain.Tag, error)
	Create(name string) (domain.Tag, error)
	Update(id int, name string) error
	Delete(id int) error
}

type CategoryRepository interface {
	FindAll() ([]domain.Category, error)
	Create(name string) (domain.Category, error)
	Update(id int, name string) error
	Delete(id int) error
}

type PostRepository interface {
	FindAll() ([]domain.Post, error)
	FindAllPublished() ([]domain.Post, error)
	FindByID(id int) (*domain.Post, error)
	FindBySlug(slug string) (*domain.Post, error)
	Create(post *domain.Post) error
	Update(post *domain.Post) error
	Delete(id int) error
}

type UserProfileRepository interface {
	FindByUserID(userID int) (*domain.UserProfile, error)
	Upsert(profile *domain.UserProfile) error
}

type TokenService interface {
	Generate(userID int, name string) (string, error)
	Validate(tokenString string) (int, error)
}
