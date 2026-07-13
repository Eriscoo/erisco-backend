package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailExists        = errors.New("email already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrNotFound           = errors.New("not found")
	ErrDuplicateEntry     = errors.New("duplicate entry")
)
