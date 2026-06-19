// Package domain also defines a single, shared error type for the whole
// service. Centralising errors here means every layer "speaks the same
// language": the domain raises an AppError, and the HTTP layer can later map
// each ErrorCode to the right status code (404, 409, 400, 500...).
package domain

import "fmt"

// ErrorCode is a typed string that classifies what kind of problem happened.
// Using a type (instead of plain strings) avoids typos and makes the set of
// codes easy to find.
type ErrorCode string

// The four categories of errors the application recognises.
const (
	ErrValidation ErrorCode = "VALIDATION_ERROR" // bad input from the caller
	ErrNotFound   ErrorCode = "NOT_FOUND"        // the resource does not exist
	ErrConflict   ErrorCode = "CONFLICT"         // a business rule was broken
	ErrInternal   ErrorCode = "INTERNAL_ERROR"   // something failed on our side
)

// AppError is our custom error. The JSON tags mean that when we send it back in
// an HTTP response it serialises to {"code": "...", "message": "..."}, giving
// the client a predictable, machine-readable shape.
type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

// Error makes AppError satisfy Go's built-in `error` interface. Any type with
// an `Error() string` method counts as an error, so our AppError can be used
// anywhere a standard error is expected.
func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// The functions below are "constructors" — small helpers so the rest of the
// code can write dmerrors.NewNotFoundError("device not found") instead of
// building the struct by hand every time. This keeps error creation short and
// consistent across the project.

// NewValidationError is used when the caller sent invalid data.
func NewValidationError(message string) *AppError {
	return &AppError{Code: ErrValidation, Message: message}
}

// NewNotFoundError is used when a requested resource does not exist.
func NewNotFoundError(message string) *AppError {
	return &AppError{Code: ErrNotFound, Message: message}
}

// NewConflictError is used when an action clashes with a business rule
// (for example, a duplicate code or an illegal status change).
func NewConflictError(message string) *AppError {
	return &AppError{Code: ErrConflict, Message: message}
}

// NewInternalError is used for unexpected failures (like a database error) that
// are not the caller's fault.
func NewInternalError(message string) *AppError {
	return &AppError{Code: ErrInternal, Message: message}
}
