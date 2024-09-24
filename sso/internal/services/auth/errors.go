package auth

import "errors"

var (
	ErrInvalidCredentials = errors.New("Invalid credentials")
	ErrUserAlreadyExists  = errors.New("User already exists")
	ErrInternal           = errors.New("Internal error")
	ErrAppNotFound        = errors.New("App not found")
)
