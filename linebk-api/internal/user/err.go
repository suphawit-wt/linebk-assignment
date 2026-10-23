package user

import "errors"

// Repository & Database Errors
var (
	ErrUserNotFound = errors.New("user not found")
)
