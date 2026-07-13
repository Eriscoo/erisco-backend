package persistence

import (
	"database/sql"
	"time"

	"github.com/eriscoo/blog-backend/internal/domain"
)

type userProfileRepository struct {
	db *sql.DB
}

func NewUserProfileRepository(db *sql.DB) *userProfileRepository {
	return &userProfileRepository{db: db}
}

func (r *userProfileRepository) FindByUserID(userID int) (*domain.UserProfile, error) {
	p := &domain.UserProfile{}
	err := r.db.QueryRow(
		"SELECT user_id, COALESCE(bio,''), COALESCE(avatar_url,''), COALESCE(website,''), COALESCE(location,''), COALESCE(phone,''), created_at, updated_at FROM user_profile WHERE user_id = $1",
		userID,
	).Scan(&p.UserID, &p.Bio, &p.AvatarURL, &p.Website, &p.Location, &p.Phone, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return p, err
}

func (r *userProfileRepository) Upsert(profile *domain.UserProfile) error {
	now := time.Now()
	_, err := r.db.Exec(
		`INSERT INTO user_profile (user_id, bio, avatar_url, website, location, phone, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 ON CONFLICT (user_id)
		 DO UPDATE SET bio = EXCLUDED.bio, avatar_url = EXCLUDED.avatar_url, website = EXCLUDED.website,
		   location = EXCLUDED.location, phone = EXCLUDED.phone, updated_at = EXCLUDED.updated_at`,
		profile.UserID, profile.Bio, profile.AvatarURL, profile.Website, profile.Location, profile.Phone, now,
	)
	return err
}
