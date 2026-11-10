package domain

import (
	"errors"
	"fmt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrDuplicateEmail     = errors.New("duplicate email")
	ErrDuplicateShortCode = errors.New("duplicate short code")
	ErrInvalidData        = errors.New("invalid user data")
	ErrDatabaseError      = errors.New("Database error")
	ErrUnexpectedError    = errors.New("unexpected error")
	ErrDeleteReferenced   = errors.New("cannot delete: it is referenced by other records")
)

type CustomError struct {
	Operation string
	Err       error
	Details   error
}

func NewCustomError(operation string, err error, details error) *CustomError {
	return &CustomError{Operation: operation, Err: err, Details: details}
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("user %s: %v, %v", e.Operation, e.Err, e.Details)
}

func (e *CustomError) Unwrap() error {
	return e.Err
}
