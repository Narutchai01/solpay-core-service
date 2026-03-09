package entities

import "errors"

// ErrorType categorizes application-level errors.
type ErrorType int

const (
	ErrTypeInternal ErrorType = iota
	ErrTypeNotFound
	ErrTypeConflict
	ErrTypeBadRequest
)

// AppError is a structured error that carries a type, message, and cause.
type AppError struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

// Unwrap returns the underlying error for use with errors.Is/As.
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new AppError with the given type, message, and cause.
func NewAppError(t ErrorType, msg string, err error) *AppError {
	return &AppError{Type: t, Message: msg, Err: err}
}

// Sentinel errors for repository-level error classification.
var (
	ErrConflict   = errors.New("record already exists")
	ErrNotFound   = errors.New("record not found")
	ErrBadRequest = errors.New("bad request")
	ErrInternal   = errors.New("internal server error")
)
