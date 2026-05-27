package domain

import "fmt"

type ErrorCode string

const (
	ErrValidation ErrorCode = "VALIDATION_ERROR"
	ErrNotFound   ErrorCode = "NOT_FOUND"
	ErrConflict   ErrorCode = "CONFLICT"
	ErrInternal   ErrorCode = "INTERNAL_ERROR"
)

type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewValidationError(message string) *AppError {
	return &AppError{Code: ErrValidation, Message: message}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{Code: ErrNotFound, Message: message}
}

func NewConflictError(message string) *AppError {
	return &AppError{Code: ErrConflict, Message: message}
}

func NewInternalError(message string) *AppError {
	return &AppError{Code: ErrInternal, Message: message}
}
