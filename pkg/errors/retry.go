package errors

// Retryable indicates whether the error is typically retriable (backoff, retry).
func Retryable(err error) bool {
	if err == nil {
		return false
	}
	var e *Error
	if !As(err, &e) {
		return false
	}
	switch e.Code {
	case CodeUnavailable, CodeTimeout, CodeInternal:
		return true
	default:
		return false
	}
}
