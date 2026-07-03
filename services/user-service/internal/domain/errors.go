package domain

import "errors"

var (
	ErrEmailRequired     = errors.New("email required")
	ErrUsernameRequired  = errors.New("username required")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)
