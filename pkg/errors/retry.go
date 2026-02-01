package errors

// Retryable indica se o erro é tipicamente retentável (backoff, retry).
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
