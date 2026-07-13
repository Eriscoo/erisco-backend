package persistence

import (
	"database/sql"
	"strings"

	"github.com/eriscoo/blog-backend/internal/domain"
)

type tagRepository struct {
	db *sql.DB
}

func NewTagRepository(db *sql.DB) *tagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) FindAll() ([]domain.Tag, error) {
	rows, err := r.db.Query("SELECT id, name FROM tags ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make([]domain.Tag, 0)
	for rows.Next() {
		var t domain.Tag
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, rows.Err()
}

func (r *tagRepository) Create(name string) (domain.Tag, error) {
	var t domain.Tag
	err := r.db.QueryRow(
		"INSERT INTO tags (name) VALUES ($1) RETURNING id, name", name,
	).Scan(&t.ID, &t.Name)
	if err != nil && strings.Contains(err.Error(), "duplicate key") {
		return t, domain.ErrDuplicateEntry
	}
	return t, err
}

func (r *tagRepository) Update(id int, name string) error {
	res, err := r.db.Exec("UPDATE tags SET name = $1 WHERE id = $2", name, id)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return domain.ErrDuplicateEntry
		}
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *tagRepository) Delete(id int) error {
	res, err := r.db.Exec("DELETE FROM tags WHERE id = $1", id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}
