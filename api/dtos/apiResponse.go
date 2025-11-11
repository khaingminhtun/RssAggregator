package dtos

// APIResponse is a standard API response structure
type APIResponse struct {
	Success bool        `json:"success"`           // true if request is successful
	Message string      `json:"message,omitempty"` // optional message
	Data    interface{} `json:"data,omitempty"`    // actual payload
	Errors  interface{} `json:"errors,omitempty"`  // optional errors
}
