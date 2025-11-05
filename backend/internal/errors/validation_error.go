package errors

import (
	"fmt"
	"net/http"
)

type ValidationError struct {
	StatusCode int
	Message    map[string][]string
}

func (err *ValidationError) Error() string {
	return fmt.Sprintf("Validation error: status_code=%d, message=%s", err.StatusCode, err.Message)
}

func NewValidationError() *ValidationError {
	return &ValidationError{
		StatusCode: http.StatusUnprocessableEntity,
		Message:    make(map[string][]string),
	}
}

func (err *ValidationError) Add(key string, errors []string) {
	existingErrors, ok := err.Message[key]
	if !ok {
		err.Message[key] = errors
		return
	}
	err.Message[key] = append(existingErrors, errors...)
}

func (err *ValidationError) HasErrors() bool {
	return len(err.Message) > 0
}
