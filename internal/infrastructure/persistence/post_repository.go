package persistence

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/eriscoo/blog-backend/internal/domain"
)

var jakartaLoc *time.Location

func init() {
	var err error
	jakartaLoc, err = time.LoadLocation("Asia/Jakarta")
	if err != nil {
		jakartaLoc = time.UTC
	}
}

func toJakarta(t time.Time) time.Time {
	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	return time.Date(year, month, day, hour, min, sec, t.Nanosecond(), jakartaLoc)
}

type postRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *postRepository {
	return &postRepository{db: db}
}

func (r *postRepository) FindAll() ([]domain.Post, error) {
	rows, err := r.db.Query(`
		SELECT p.id, p.title, p.slug, COALESCE(p.body, ''), COALESCE(p.image_url, ''),
		       COALESCE(p.categories, ''), COALESCE(p.tags, ''), p.created_by, u.name,
		       COALESCE(up.avatar_url, ''), p.status,
		       p.published_at, p.created_at, p.updated_at
		FROM posts p
		JOIN users u ON u.id = p.created_by
		LEFT JOIN user_profile up ON up.user_id = p.created_by
		ORDER BY p.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]domain.Post, 0)
	for rows.Next() {
		var p domain.Post
		if err := rows.Scan(
			&p.ID, &p.Title, &p.Slug, &p.Body, &p.ImageURL,
			&p.Categories, &p.Tags, &p.CreatedBy, &p.CreatedByName,
			&p.AuthorAvatarURL,
			&p.Status,
			&p.PublishedAt, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		p.CreatedAt = toJakarta(p.CreatedAt)
		p.UpdatedAt = toJakarta(p.UpdatedAt)
		if p.PublishedAt != nil {
			t := toJakarta(*p.PublishedAt)
			p.PublishedAt = &t
		}
		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := r.populateCategoryNames(posts); err != nil {
		return nil, err
	}
	if err := r.populateTagNames(posts); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *postRepository) FindAllPublished() ([]domain.Post, error) {
	rows, err := r.db.Query(`
		SELECT p.id, p.title, p.slug, COALESCE(p.body, ''), COALESCE(p.image_url, ''),
		       COALESCE(p.categories, ''), COALESCE(p.tags, ''), p.created_by, u.name,
		       COALESCE(up.avatar_url, ''), p.status,
		       p.published_at, p.created_at, p.updated_at
		FROM posts p
		JOIN users u ON u.id = p.created_by
		LEFT JOIN user_profile up ON up.user_id = p.created_by
		WHERE p.status = 'published'
		ORDER BY p.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]domain.Post, 0)
	for rows.Next() {
		var p domain.Post
		if err := rows.Scan(
			&p.ID, &p.Title, &p.Slug, &p.Body, &p.ImageURL,
			&p.Categories, &p.Tags, &p.CreatedBy, &p.CreatedByName,
			&p.AuthorAvatarURL,
			&p.Status,
			&p.PublishedAt, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		p.CreatedAt = toJakarta(p.CreatedAt)
		p.UpdatedAt = toJakarta(p.UpdatedAt)
		if p.PublishedAt != nil {
			t := toJakarta(*p.PublishedAt)
			p.PublishedAt = &t
		}
		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := r.populateCategoryNames(posts); err != nil {
		return nil, err
	}
	if err := r.populateTagNames(posts); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *postRepository) populateCategoryNames(posts []domain.Post) error {
	catRows, err := r.db.Query("SELECT id, name FROM categories")
	if err != nil {
		return err
	}
	defer catRows.Close()

	catMap := make(map[int]string)
	for catRows.Next() {
		var id int
		var name string
		if err := catRows.Scan(&id, &name); err != nil {
			return err
		}
		catMap[id] = name
	}
	if err := catRows.Err(); err != nil {
		return err
	}

	for i, p := range posts {
		if p.Categories == "" {
			continue
		}
		ids := strings.Split(p.Categories, ",")
		names := make([]string, 0, len(ids))
		for _, idStr := range ids {
			idStr = strings.TrimSpace(idStr)
			var id int
			if _, err := fmt.Sscanf(idStr, "%d", &id); err == nil {
				if name, ok := catMap[id]; ok {
					names = append(names, name)
				}
			}
		}
		posts[i].CategoryNames = strings.Join(names, ", ")
	}

	return nil
}

func (r *postRepository) populateTagNames(posts []domain.Post) error {
	tagRows, err := r.db.Query("SELECT id, name FROM tags")
	if err != nil {
		return err
	}
	defer tagRows.Close()

	tagMap := make(map[int]string)
	for tagRows.Next() {
		var id int
		var name string
		if err := tagRows.Scan(&id, &name); err != nil {
			return err
		}
		tagMap[id] = name
	}
	if err := tagRows.Err(); err != nil {
		return err
	}

	for i, p := range posts {
		if p.Tags == "" {
			continue
		}
		ids := strings.Split(p.Tags, ",")
		names := make([]string, 0, len(ids))
		for _, idStr := range ids {
			idStr = strings.TrimSpace(idStr)
			var id int
			if _, err := fmt.Sscanf(idStr, "%d", &id); err == nil {
				if name, ok := tagMap[id]; ok {
					names = append(names, name)
				}
			}
		}
		posts[i].TagNames = strings.Join(names, ", ")
	}

	return nil
}

func (r *postRepository) FindByID(id int) (*domain.Post, error) {
	var p domain.Post
	err := r.db.QueryRow(`
		SELECT p.id, p.title, p.slug, COALESCE(p.body, ''), COALESCE(p.image_url, ''),
		       COALESCE(p.categories, ''), COALESCE(p.tags, ''), p.created_by, u.name,
		       COALESCE(up.avatar_url, ''), p.status,
		       p.published_at, p.created_at, p.updated_at
		FROM posts p
		JOIN users u ON u.id = p.created_by
		LEFT JOIN user_profile up ON up.user_id = p.created_by
		WHERE p.id = $1
	`, id).Scan(
		&p.ID, &p.Title, &p.Slug, &p.Body, &p.ImageURL,
		&p.Categories, &p.Tags, &p.CreatedBy, &p.CreatedByName,
		&p.AuthorAvatarURL,
		&p.Status,
		&p.PublishedAt, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	p.CreatedAt = toJakarta(p.CreatedAt)
	p.UpdatedAt = toJakarta(p.UpdatedAt)
	if p.PublishedAt != nil {
		t := toJakarta(*p.PublishedAt)
		p.PublishedAt = &t
	}
	return &p, nil
}

func (r *postRepository) FindBySlug(slug string) (*domain.Post, error) {
	var p domain.Post
	err := r.db.QueryRow(`
		SELECT p.id, p.title, p.slug, COALESCE(p.body, ''), COALESCE(p.image_url, ''),
		       COALESCE(p.categories, ''), COALESCE(p.tags, ''), p.created_by, u.name,
		       COALESCE(up.avatar_url, ''), p.status,
		       p.published_at, p.created_at, p.updated_at
		FROM posts p
		JOIN users u ON u.id = p.created_by
		LEFT JOIN user_profile up ON up.user_id = p.created_by
		WHERE p.slug = $1 AND p.status = 'published'
	`, slug).Scan(
		&p.ID, &p.Title, &p.Slug, &p.Body, &p.ImageURL,
		&p.Categories, &p.Tags, &p.CreatedBy, &p.CreatedByName,
		&p.AuthorAvatarURL,
		&p.Status,
		&p.PublishedAt, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	p.CreatedAt = toJakarta(p.CreatedAt)
	p.UpdatedAt = toJakarta(p.UpdatedAt)
	if p.PublishedAt != nil {
		t := toJakarta(*p.PublishedAt)
		p.PublishedAt = &t
	}

	single := []domain.Post{p}
	if err := r.populateCategoryNames(single); err != nil {
		return nil, err
	}
	if err := r.populateTagNames(single); err != nil {
		return nil, err
	}
	p.CategoryNames = single[0].CategoryNames
	p.TagNames = single[0].TagNames

	return &p, nil
}

func (r *postRepository) Create(post *domain.Post) error {
	err := r.db.QueryRow(`
		INSERT INTO posts (title, slug, body, image_url, categories, tags, created_by, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`, post.Title, post.Slug, post.Body, post.ImageURL,
		post.Categories, post.Tags, post.CreatedBy, post.Status,
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil && strings.Contains(err.Error(), "duplicate key") {
		return domain.ErrDuplicateEntry
	}
	return err
}

func (r *postRepository) Update(post *domain.Post) error {
	res, err := r.db.Exec(`
		UPDATE posts SET title = $1, slug = $2, body = $3, image_url = $4, categories = $5,
		                 tags = $6, status = $7, updated_at = NOW()
		WHERE id = $8
	`, post.Title, post.Slug, post.Body, post.ImageURL, post.Categories,
		post.Tags, post.Status, post.ID,
	)
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

func (r *postRepository) Delete(id int) error {
	res, err := r.db.Exec("DELETE FROM posts WHERE id = $1", id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}
