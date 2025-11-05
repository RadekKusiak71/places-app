package utils

import (
	"encoding/json"
	"net/http"
)

func ReadJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func WriteJSON(w http.ResponseWriter, statusCode int, v any) error {
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(v)
}
