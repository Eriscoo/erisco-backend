package auth

import (
	"github.com/eriscoo/blog-backend/internal/application"
	"github.com/eriscoo/blog-backend/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepo application.UserRepository
	tokens   application.TokenService
}

func New(repo application.UserRepository, tokens application.TokenService) *Service {
	return &Service{
		userRepo: repo,
		tokens:   tokens,
	}
}

func (s *Service) Register(name, email, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := &domain.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		RoleID:       2,
	}

	if err := s.userRepo.Create(user); err != nil {
		return "", err
	}

	return s.tokens.Generate(user.ID, user.Name)
}

func (s *Service) Login(email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", domain.ErrInvalidCredentials
	}

	return s.tokens.Generate(user.ID, user.Name)
}

func (s *Service) GetUser(id int) (*domain.User, error) {
	return s.userRepo.FindByID(id)
}
