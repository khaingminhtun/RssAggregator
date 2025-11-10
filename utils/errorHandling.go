package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

// RespondWithError sends JSON error response for clients
func RespondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := ErrorResponse{Error: message}
	json.NewEncoder(w).Encode(resp)
}

// LogAndRespond logs server errors and responds with generic message
func LogAndRespond(w http.ResponseWriter, err error, message string) {
	log.Printf("Server error: %v\n", err) // log full error
	RespondWithError(w, http.StatusInternalServerError, message)
}
