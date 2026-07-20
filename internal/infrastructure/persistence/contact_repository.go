package persistence

import (
	"database/sql"

	"github.com/eriscoo/blog-backend/internal/domain"
)

type contactRepository struct {
	db *sql.DB
}

func NewContactRepository(db *sql.DB) *contactRepository {
	return &contactRepository{db: db}
}

func (r *contactRepository) Create(msg *domain.ContactMessage) error {
	return r.db.QueryRow(`
		INSERT INTO contact_messages (name, email, subject, phone, message)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`, msg.Name, msg.Email, msg.Subject, msg.Phone, msg.Message,
	).Scan(&msg.ID, &msg.CreatedAt)
}
