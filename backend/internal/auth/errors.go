package auth

import (
	"net/http"

	"github.com/RadekKusiak71/places-app/internal/errors"
)

func ErrUsernameIsTaken() errors.APIError {
	return errors.APIError{
		StatusCode: http.StatusConflict,
		Message:    "username is already taken",
	}
}
