package errorHandle

import (
	"net/http"

	"github.com/khaingminhtun/rssagg/internal/json"
	"github.com/khaingminhtun/rssagg/internal/log"
)

const (
	TypeInvalidRequest = "invalid_request"
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

var typeToStatus = map[string]int{
	TypeInvalidRequest: http.StatusUnprocessableEntity,
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
