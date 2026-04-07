package domain

import (
	"errors"
	"fmt"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrDuplicateEmail   = errors.New("duplicate email")
	ErrInvalidUserData  = errors.New("invalid user data")
	ErrDatabaseError    = errors.New("Database error")
	ErrUnexpectedError  = errors.New("unexpected error")
	ErrDeleteReferenced = errors.New("cannot delete: it is referenced by other records")
)

type UserError struct {
	Operation string
	Err       error
	Details   error
}

func NewUserError(operation string, err error, details error) *UserError {
	return &UserError{Operation: operation, Err: err, Details: details}
}

func (e *UserError) Error() string {
	return fmt.Sprintf("user %s: %v, %v", e.Operation, e.Err, e.Details)
}

func (e *UserError) Unwrap() error {
	return e.Err
}
