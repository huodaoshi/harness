package store

import "errors"

var (
	ErrSessionNotFound    = errors.New("session not found")
	ErrSessionForbidden   = errors.New("session does not belong to user")
	ErrSessionEnded       = errors.New("session already ended")
	ErrSessionMessageCap  = errors.New("session message cap exceeded")
)
