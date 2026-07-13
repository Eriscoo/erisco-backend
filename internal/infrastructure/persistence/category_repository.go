package persistence

import (
	"database/sql"
	"strings"

	"github.com/eriscoo/blog-backend/internal/domain"
)

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *categoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) FindAll() ([]domain.Category, error) {
	rows, err := r.db.Query("SELECT id, name FROM categories ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]domain.Category, 0)
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, rows.Err()
}

func (r *categoryRepository) Create(name string) (domain.Category, error) {
	var c domain.Category
	err := r.db.QueryRow(
		"INSERT INTO categories (name) VALUES ($1) RETURNING id, name", name,
	).Scan(&c.ID, &c.Name)
	if err != nil && strings.Contains(err.Error(), "duplicate key") {
		return c, domain.ErrDuplicateEntry
	}
	return c, err
}

func (r *categoryRepository) Update(id int, name string) error {
	res, err := r.db.Exec("UPDATE categories SET name = $1 WHERE id = $2", name, id)
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

func (r *categoryRepository) Delete(id int) error {
	res, err := r.db.Exec("DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}
