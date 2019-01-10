package domain

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	// ErrNotFound is an error when the requested resource is not found
	ErrNotFound = errors.New("The requested resource is not found")
)

// respondwithJSON write json response format
func RespondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// respondwithError return error message
func RespondWithError(w http.ResponseWriter, code int, msg string) {
	RespondwithJSON(w, code, map[string]string{"message": msg})
}
