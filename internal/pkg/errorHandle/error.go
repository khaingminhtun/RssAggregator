package errorHandle

import (
	"net/http"

	"github.com/khaingminhtun/rssagg/internal/pkg/json"
	"github.com/khaingminhtun/rssagg/internal/pkg/log"
)

const (
	TypeInvalidRequest       = "invalid_request"
	TypeNotFound             = "not_found"
	TypeConflict             = "conflict"
	TypeUnauthorized         = "unauthorized"
	TypeFeedAlreadyExists    = "feed_already_exists"
	TypeFeedDiscoveryFailed  = "feed_discovery_failed"
	TypeFeedParseFailed      = "feed_parse_failed"
	TypeDatabaseError        = "database_error"
	TypeExternalServiceError = "external_service_error"
)

type ApiError struct {
	Type string
	Msg  string
}

func (e ApiError) Error() string {
	return e.Msg
}

func BadRequest(msg string) error {
	return ApiError{
		Type: TypeInvalidRequest,
		Msg:  msg,
	}
}

func NotFound(msg string) error {
	return ApiError{Type: TypeNotFound, Msg: msg}
}

func Conflict(msg string) error {
	return ApiError{Type: TypeConflict, Msg: msg}
}

func Unauthorized(msg string) error {
	return ApiError{Type: TypeUnauthorized, Msg: msg}
}

// ----------------- Custom Feed Errors -----------------
func FeedAlreadyExists(msg string) error {
	return ApiError{Type: TypeFeedAlreadyExists, Msg: msg}
}

func FeedDiscoveryFailed(msg string) error {
	return ApiError{Type: TypeFeedDiscoveryFailed, Msg: msg}
}

func FeedParseFailed(msg string) error {
	return ApiError{Type: TypeFeedParseFailed, Msg: msg}
}

func DatabaseError(msg string) error {
	return ApiError{Type: TypeDatabaseError, Msg: msg}
}

func ExternalServiceError(msg string) error {
	return ApiError{Type: TypeExternalServiceError, Msg: msg}
}

var typeToStatus = map[string]int{
	TypeInvalidRequest:       http.StatusUnprocessableEntity, // 422
	TypeNotFound:             http.StatusNotFound,            // 404
	TypeConflict:             http.StatusConflict,            // 409
	TypeUnauthorized:         http.StatusUnauthorized,        // 401
	TypeFeedAlreadyExists:    http.StatusConflict,            // 409
	TypeFeedDiscoveryFailed:  http.StatusBadRequest,          // 400
	TypeFeedParseFailed:      http.StatusBadRequest,          // 400
	TypeDatabaseError:        http.StatusInternalServerError, // 500
	TypeExternalServiceError: http.StatusBadGateway,          // 502
}

func RespondHTTPsError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	if apiErr, ok := err.(ApiError); ok {
		status, _ := typeToStatus[apiErr.Type]
		json.RespondError(w, status, apiErr)
		return
	}

	// fallback for unknown errors
	json.RespondError(w, http.StatusInternalServerError, "internal server error")
	// log the error with route
	log.Error("internal server error",
		"error", err,
		"method", r.Method,
		"path", r.URL.Path,
	)
}
