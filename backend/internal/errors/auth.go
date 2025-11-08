package errors

import "net/http"

func InvalidCredentialsError() APIError {
	return APIError{
		StatusCode: http.StatusUnauthorized,
		Message:    "invalid username or password",
	}
}

func InvalidTokenError() APIError {
	return APIError{
		StatusCode: http.StatusUnauthorized,
		Message:    "invalid token",
	}
}

func InvalidTokenErrorWithMessage(err error) APIError {
	return APIError{
		StatusCode: http.StatusUnauthorized,
		Message:    err.Error(),
	}
}

func MissingAuthorizationHeader() APIError {
	return APIError{
		StatusCode: http.StatusUnauthorized,
		Message:    "missing authorization header",
	}
}

func InvalidAuthorizationHeader() APIError {
	return APIError{
		StatusCode: http.StatusUnauthorized,
		Message:    "invalid authorization header expected: Bearer <token>",
	}
}
