package exceptions

import "errors"

var (
	ErrEventNotFound      = errors.New("event not found")
	ErrEventFull          = errors.New("event has reached its maximum capacity")
	ErrAlreadyJoined      = errors.New("user has already joined this event")
	ErrUnauthorizedAction = errors.New("unauthorized action: only admins can perform this")
)
