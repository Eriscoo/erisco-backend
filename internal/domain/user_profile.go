package domain

import "time"

type UserProfile struct {
	UserID    int       `json:"user_id"`
	Bio       string    `json:"bio"`
	AvatarURL string    `json:"avatar_url"`
	Website   string    `json:"website"`
	Location  string    `json:"location"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
