package domain

import "time"

type Post struct {
	ID            int        `json:"id"`
	Title         string     `json:"title"`
	Slug          string     `json:"slug"`
	Body          string     `json:"body"`
	ImageURL      string     `json:"image_url"`
	Categories    string     `json:"categories"`
	CategoryNames string     `json:"category_names"`
	Tags          string     `json:"tags"`
	TagNames      string     `json:"tag_names"`
	CreatedBy      int        `json:"created_by"`
	CreatedByName  string     `json:"created_by_name"`
	AuthorAvatarURL string    `json:"author_avatar_url"`
	Status        string     `json:"status"`
	PublishedAt   *time.Time `json:"-"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
