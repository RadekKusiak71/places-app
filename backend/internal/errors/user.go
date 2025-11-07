package errors

import (
	"net/http"
)

func ErrUserAlreadyExists() error {
	return APIError{
		StatusCode: http.StatusConflict,
		Message:    "user with the given username already exists",
	}
}
