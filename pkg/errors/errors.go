// Package errors provides typed errors, Is/As, mapping for HTTP/exit code/retry
// and optional stack. Avoids scattered errors.New.
package errors

import (
	stderrors "errors"
	"fmt"
	"runtime"
)

var _ error = (*Error)(nil)

// Error is a typed error with code, message and cause.
type Error struct {
	Code    string
	Message string
	Cause   error
	Stack   []byte
}

// Error implements error.
func (e *Error) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

// Unwrap allows errors.Unwrap and errors.Is/As (stdlib).
func (e *Error) Unwrap() error { return e.Cause }

// WithStack records the stack trace on the error (optional).
func (e *Error) WithStack() *Error {
	e.Stack = stack()
	return e
}

func stack() []byte {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return buf[:n]
}

// New creates a typed error. Optional: chain with Wrap(err) and .WithStack().
func New(code, message string) *Error {
	return &Error{Code: code, Message: message}
}

// Wrap wraps an error with code and message.
func Wrap(err error, code, message string) *Error {
	if err == nil {
		return nil
	}
	return &Error{Code: code, Message: message, Cause: err}
}

// Wrapf wraps with a formatted message.
func Wrapf(err error, code, format string, args ...any) *Error {
	if err == nil {
		return nil
	}
	return Wrap(err, code, fmt.Sprintf(format, args...))
}

// Is delegates to errors.Is (stdlib).
func Is(err, target error) bool { return stderrors.Is(err, target) }

// As delegates to errors.As (stdlib).
func As(err error, target any) bool { return stderrors.As(err, target) }

// Common codes for HTTP/exit mapping.
const (
	CodeNotFound     = "NOT_FOUND"
	CodeInvalidInput = "INVALID_INPUT"
	CodeUnauthorized = "UNAUTHORIZED"
	CodeForbidden    = "FORBIDDEN"
	CodeConflict     = "CONFLICT"
	CodeInternal     = "INTERNAL"
	CodeUnavailable  = "UNAVAILABLE"
	CodeTimeout      = "TIMEOUT"
)

// Sentinels for use with errors.Is.
var (
	ErrNotFound     = New(CodeNotFound, "not found")
	ErrInvalidInput = New(CodeInvalidInput, "invalid input")
	ErrUnauthorized = New(CodeUnauthorized, "unauthorized")
	ErrForbidden    = New(CodeForbidden, "forbidden")
	ErrConflict     = New(CodeConflict, "conflict")
	ErrInternal     = New(CodeInternal, "internal error")
	ErrUnavailable  = New(CodeUnavailable, "unavailable")
	ErrTimeout      = New(CodeTimeout, "timeout")
)
