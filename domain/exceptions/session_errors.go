package exceptions

import "errors"

var (
	ErrInvalidToken   = errors.New("invalid or expired token")
	ErrTokenExpired   = errors.New("token has expired")
	ErrSessionExpired = errors.New("session has expired, please login again")
)
