package auth

import (
	"net/http"

	"github.com/khaingminhtun/rssagg/internal/errorHandle"
	"github.com/khaingminhtun/rssagg/internal/json"
)

type AuthHandler struct {
	AuthService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
	}
}

//register user
func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Implementation of the handler to register a user
	var req RegisterRequest

	if !json.DecodeJSON(w, r, &req) {
		return
	}

	err := h.AuthService.RegisterUser(r.Context(), req)

	if err != nil {
		switch err {
		case errorHandle.ErrInvalidInput:
			json.RespondError(w, http.StatusBadRequest, "invalid input")
		case errorHandle.ErrUserAlreadyExists:
			json.RespondError(w, http.StatusConflict, "user already exists")
		default:
			json.RespondError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	json.RespondJSON(w, http.StatusCreated, json.SuccessResponse[any]{Data: "user registered successfully"})

}
