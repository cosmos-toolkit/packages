package errors

// ExitCode maps error to exit code (CLI).
// 0 = success, 1 = generic error, 2 = invalid input, 3 = not found, etc.
func ExitCode(err error) int {
	if err == nil {
		return 0
	}
	var e *Error
	if !As(err, &e) {
		return 1
	}
	switch e.Code {
	case CodeInvalidInput:
		return 2
	case CodeNotFound:
		return 3
	case CodeUnauthorized, CodeForbidden:
		return 4
	case CodeConflict:
		return 5
	case CodeUnavailable, CodeTimeout:
		return 6
	default:
		return 1
	}
}
