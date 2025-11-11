package utils

import (
	"encoding/json"
	"net/http"
)

// RespondWithJSON writes any payload as JSON with the given HTTP status code
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// Set content-type header first
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	// Encode payload to JSON
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		// If encoding fails, return internal server error
		http.Error(w, `{"success":false,"message":"failed to encode JSON"}`, http.StatusInternalServerError)
	}
}
