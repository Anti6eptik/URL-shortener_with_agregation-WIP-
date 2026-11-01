package service

import (
	"fmt"
	"net/http"
)

type Kind string

const (
	KindBadRequest      Kind = "bad_request"
	KindUnauthorized    Kind = "unauthorized"
	KindForbidden       Kind = "forbidden"
	KindNotFound        Kind = "not_found"
	KindConflict        Kind = "conflict"
	KindUnprocessable   Kind = "unprocessable"
	KindTooManyRequests Kind = "too_many_requests"
	KindInternal        Kind = "internal"
)

type ServiceError struct {
	Kind    Kind   `json:"kind"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *ServiceError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Kind, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Kind, e.Message)
}

func (e *ServiceError) Unwrap() error {
	return e.Err
}

func (e *ServiceError) HTTPStatus() int {
	switch e.Kind {
	case KindBadRequest:
		return http.StatusBadRequest
	case KindUnauthorized:
		return http.StatusUnauthorized
	case KindForbidden:
		return http.StatusForbidden
	case KindNotFound:
		return http.StatusNotFound
	case KindConflict:
		return http.StatusConflict
	case KindUnprocessable:
		return http.StatusUnprocessableEntity
	case KindTooManyRequests:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}

func NewBadRequestError(msg string, err error) *ServiceError {
	return &ServiceError{Kind: KindBadRequest, Message: msg, Err: err}
}

func NewUnauthorizedError(msg string, err error) *ServiceError {
	return &ServiceError{Kind: KindUnauthorized, Message: msg, Err: err}
}

func NewForbiddenError(msg string, err error) *ServiceError {
	return &ServiceError{Kind: KindForbidden, Message: msg, Err: err}
}

func NewNotFoundError(msg string, err error) *ServiceError {
	return &ServiceError{Kind: KindNotFound, Message: msg, Err: err}
}

func NewConflictError(msg string, err error) *ServiceError {
	return &ServiceError{Kind: KindConflict, Message: msg, Err: err}
}

func NewUnprocessableError(msg string, err error) *ServiceError {
	return &ServiceError{Kind: KindUnprocessable, Message: msg, Err: err}
}

func NewInternalError(msg string, err error) *ServiceError {
	return &ServiceError{Kind: KindInternal, Message: msg, Err: err}
}
