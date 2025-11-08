package errors

import "net/http"

func PlaceNotFoundError() APIError {
	return APIError{
		StatusCode: http.StatusNotFound,
		Message:    "place not found",
	}
}
