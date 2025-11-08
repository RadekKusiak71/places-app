package errors

import (
	"fmt"
	"net/http"
)

type APIError struct {
	StatusCode int `json:"status_code"`
	Message    any `json:"message"`
}

func (err APIError) Error() string {
	return fmt.Sprintf("API error: status_code=%d, message=%s", err.StatusCode, err.Message)
}

func NewAPIError(statusCode int, message any) error {
	return APIError{StatusCode: statusCode, Message: message}
}

func InternalServerError() error {
	return NewAPIError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func InvalidRequestError() error {
	return NewAPIError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
}
