package apperrors

import (
	"errors"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type Error struct {
	Status   int
	Messages []string
}

func (e *Error) Error() string {
	return strings.Join(e.Messages, "; ")
}

func New(status int, messages ...string) *Error {
	return &Error{Status: status, Messages: messages}
}

var (
	ErrEmailTaken    = New(http.StatusUnprocessableEntity, "email already taken")
	ErrUsernameTaken = New(http.StatusUnprocessableEntity, "username already taken")
	ErrInvalidLogin  = New(http.StatusUnauthorized, "invalid email or password")
	ErrUnauthorized  = New(http.StatusUnauthorized, "unauthorized")
	ErrForbidden     = New(http.StatusForbidden, "forbidden")
	ErrNotFound      = New(http.StatusNotFound, "not found")
	ErrInvalidBody   = New(http.StatusUnprocessableEntity, "invalid json body")
)

func Is(err, target error) bool { return errors.Is(err, target) }

func IsNotFound(err error) bool { return errors.Is(err, gorm.ErrRecordNotFound) }

func As(err error) (*Error, bool) {
	var e *Error
	if errors.As(err, &e) {
		return e, true
	}
	return nil, false
}
