package exceptions

import "errors"

var (
	ErrDuplicateEmail = errors.New("email already exists")
	ErrUserNotFound   = errors.New("user not found")
)
