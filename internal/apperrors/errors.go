package apperrors

import "errors"

var (
	ErrEmailTaken      = errors.New("email already taken")
	ErrUsernameTaken   = errors.New("username already taken")
	ErrInvalidLogin    = errors.New("invalid email or password")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden")
	ErrNotFound        = errors.New("not found")
)

func Is(err, target error) bool { return errors.Is(err, target) }
