package persistence

import (
	"database/sql"
	"strings"

	"github.com/eriscoo/blog-backend/internal/domain"
	_ "github.com/lib/pq"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	err := r.db.QueryRow(
		"INSERT INTO users (name, email, password_hash, role_id) VALUES ($1, $2, $3, $4) RETURNING id",
		user.Name, user.Email, user.PasswordHash, user.RoleID,
	).Scan(&user.ID)
	if err != nil && strings.Contains(err.Error(), "duplicate key") {
		return domain.ErrEmailExists
	}
	return err
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.QueryRow(
		"SELECT id, name, email, password_hash, role_id, created_at, updated_at FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.RoleID, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	return user, err
}

func (r *userRepository) FindByID(id int) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.QueryRow(
		"SELECT id, name, email, password_hash, role_id, created_at, updated_at FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.RoleID, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	return user, err
}

func OpenDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
