package exceptions

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrWeakPassword       = errors.New("password must be at least 8 characters")
	ErrPasswordMismatch   = errors.New("passwords do not match")
)
