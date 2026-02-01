// Package errors fornece erros tipados, Is/As, mapeamento para HTTP/exit code/retry
// e stack opcional. Evita errors.New espalhado.
package errors

import (
	stderrors "errors"
	"fmt"
	"runtime"
)

var _ error = (*Error)(nil)

// Error é um erro tipado com código, mensagem e causa.
type Error struct {
	Code    string
	Message string
	Cause   error
	Stack   []byte
}

// Error implementa error.
func (e *Error) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

// Unwrap permite errors.Unwrap e errors.Is/As (stdlib).
func (e *Error) Unwrap() error { return e.Cause }

// WithStack grava o stack trace no erro (opcional).
func (e *Error) WithStack() *Error {
	e.Stack = stack()
	return e
}

func stack() []byte {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return buf[:n]
}

// New cria um erro tipado. Opcional: encadear com Wrap(err) e .WithStack().
func New(code, message string) *Error {
	return &Error{Code: code, Message: message}
}

// Wrap envolve um erro com código e mensagem.
func Wrap(err error, code, message string) *Error {
	if err == nil {
		return nil
	}
	return &Error{Code: code, Message: message, Cause: err}
}

// Wrapf envolve com mensagem formatada.
func Wrapf(err error, code, format string, args ...any) *Error {
	if err == nil {
		return nil
	}
	return Wrap(err, code, fmt.Sprintf(format, args...))
}

// Is delega para errors.Is (stdlib).
func Is(err, target error) bool { return stderrors.Is(err, target) }

// As delega para errors.As (stdlib).
func As(err error, target any) bool { return stderrors.As(err, target) }

// Códigos comuns para mapeamento HTTP/exit.
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

// Sentinels para uso com errors.Is.
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
