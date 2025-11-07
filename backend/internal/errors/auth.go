package errors

import "net/http"

func ErrInvalidCredentials() error {
	return APIError{
		StatusCode: http.StatusUnauthorized,
		Message:    "invalid username or password",
	}
}

func InvalidTokenError() error {
	return APIError{
		StatusCode: http.StatusUnauthorized,
		Message:    "invalid token",
	}
}

func InvalidTokenErrorWithMessage(err error) error {
	return APIError{
		StatusCode: http.StatusUnauthorized,
		Message:    err.Error(),
	}
}
