package json

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse[T any] struct {
	Data T `json:"data"`
}

// RespondJSON writes a JSON response with a status code
func RespondJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		// Fallback: encoding error
		http.Error(
			w,
			`{"error":"failed to encode JSON"}`,
			http.StatusInternalServerError,
		)
	}
}

// RespondError writes a JSON error response
func RespondError(w http.ResponseWriter, code int, message string) {
	RespondJSON(w, code, ErrorResponse{
		Error: message,
	})
}

// DecodeJSON decodes request body into dst and handles common errors
func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		RespondError(w, http.StatusBadRequest, "invalid JSON body")
		return false
	}
	return true
}
