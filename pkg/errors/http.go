package errors

import "net/http"

// HTTPStatus maps error code to HTTP status.
func HTTPStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}
	var e *Error
	if !As(err, &e) {
		return http.StatusInternalServerError
	}
	switch e.Code {
	case CodeNotFound:
		return http.StatusNotFound
	case CodeInvalidInput:
		return http.StatusBadRequest
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeConflict:
		return http.StatusConflict
	case CodeUnavailable:
		return http.StatusServiceUnavailable
	case CodeTimeout:
		return http.StatusGatewayTimeout
	case CodeInternal:
	default:
	}
	return http.StatusInternalServerError
}
