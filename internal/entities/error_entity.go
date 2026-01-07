package entities

import "errors"

type ErrorType int

const (
	ErrTypeInternal ErrorType = iota
	ErrTypeNotFound
	ErrTypeConflict
	ErrTypeBadRequest
)

type AppError struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(t ErrorType, msg string, err error) *AppError {
	return &AppError{Type: t, Message: msg, Err: err}
}

var ErrConflict = errors.New("record already exists")
var ErrNotFound = errors.New("record not found")
var ErrBadRequest = errors.New("bad request")
var ErrInternal = errors.New("internal server error")
