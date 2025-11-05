package utils

import (
	"log"
	"net/http"

	"github.com/RadekKusiak71/places-app/internal/errors"
)

type APIFunc func(http.ResponseWriter, *http.Request) error

func MakeHandlerFunc(next APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := next(w, r)
		if err == nil {
			return
		}

		apiErr, ok := err.(errors.APIError)
		if ok {
			WriteJSON(w, apiErr.StatusCode, apiErr)
			return
		}

		validationError, ok := err.(*errors.ValidationError)
		if ok {
			WriteJSON(w, validationError.StatusCode, validationError)
			return
		}

		log.Printf("Not handled error occured: %s", err.Error())
		WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"status_code": http.StatusInternalServerError,
			"message":     http.StatusText(http.StatusInternalServerError),
		})
	}
}
