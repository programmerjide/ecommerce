package service

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrValidationFailed   = errors.New("validation failed")
	ErrDuplicateEmail     = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized access")
)
