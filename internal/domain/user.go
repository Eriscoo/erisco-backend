package domain

import "time"

type User struct {
	ID           int
	Name         string
	Email        string
	PasswordHash string
	RoleID       int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserRole struct {
	ID        int
	RoleName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
