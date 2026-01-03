package json

import (
	"encoding/json"
	"net/http"
)

// StandardResponse is the unified envelope for every API call
type StandardResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}

// RespondJSON writes the unified response format
func RespondJSON(w http.ResponseWriter, code int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	// Determine success based on HTTP status code
	isSuccess := code >= 200 && code < 300

	var response any
	if isSuccess {
		response = StandardResponse[any]{
			Success: true,
			Message: message,
			Data:    data,
		}
	} else {
		response = StandardResponse[any]{
			Success: false,
			Message: message,
			Error:   data, // In case of error, 'data' usually contains error details
		}
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"success":false,"message":"failed to encode JSON"}`, http.StatusInternalServerError)
	}
}

// RespondError simplifies sending error messages
func RespondError(w http.ResponseWriter, code int, message string) {
	RespondJSON(w, code, message, nil)
}

// DecodeJSON remains the same, but uses the new RespondError
func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		RespondError(w, http.StatusBadRequest, "invalid JSON body")
		return false
	}
	return true
}