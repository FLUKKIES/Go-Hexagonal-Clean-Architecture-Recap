package exceptions

import "errors"

var (
	ErrOAuthEmailConflict = errors.New("email already registered with a different method")
)
