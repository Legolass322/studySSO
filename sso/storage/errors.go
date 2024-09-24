package storage

import "errors"

var (
	ErrUserNotFound      = errors.New("User not found")
	ErrUserAlreadyExists = errors.New("User already exists")
	ErrAppNotFound       = errors.New("App not found")
)
