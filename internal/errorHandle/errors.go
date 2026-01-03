package errorHandle

import (
	"errors"
	"net/http"

	"github.com/jackc/pgconn"
	"github.com/khaingminhtun/rssagg/internal/json"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidInput       = errors.New("invalid input")
	ErrInternal           = errors.New("internal error")
	ErrUserNotFound       = errors.New("user not found")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
)

func IsUniqueViolation(err error) bool {
	pgErr, ok := err.(*pgconn.PgError)
	return ok && pgErr.Code == "23505"
}

func RespondHTTPError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidInput):
		json.RespondError(w, http.StatusBadRequest, "invalid input")

	case errors.Is(err, ErrUserAlreadyExists):
		json.RespondError(w, http.StatusConflict, "user already exists")

	case errors.Is(err, ErrUserNotFound), errors.Is(err, ErrInvalidCredentials):
		json.RespondError(w, http.StatusUnauthorized, "invalid email or password")

	case errors.Is(err, ErrUnauthorized):
		json.RespondError(w, http.StatusUnauthorized, "unauthorized")

	case errors.Is(err, ErrForbidden):
		json.RespondError(w, http.StatusForbidden, "forbidden")

	default:
		json.RespondError(w, http.StatusInternalServerError, "internal server error")
	}
}
