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

// register user
func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Implementation of the handler to register a user
	var req RegisterRequest

	if !json.DecodeJSON(w, r, &req) {
		return
	}

	err := h.AuthService.RegisterUser(r.Context(), req)
	if err != nil {
		errorHandle.RespondHTTPError(w, err)
		return
	}

	json.RespondJSON(w, http.StatusCreated, json.SuccessResponse[any]{Data: "user registered successfully"})

}

// authenicate
func (h *AuthHandler) Authenticate(w http.ResponseWriter, r *http.Request) {

	var req AuthRequest

	if !json.DecodeJSON(w, r, &req) {
		return
	}

	//authenicat and generate tokens
	resp, err := h.AuthService.AuthenticateUser(r.Context(), req)
	if err != nil {
		errorHandle.RespondHTTPError(w, err)
		return
	}

	//set refreh token in http-only secure cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    resp.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,                // in production to true
		SameSite: http.SameSiteLaxMode, // default is strict ,
		MaxAge:   7 * 24 * 60 * 60,
	})

	json.RespondJSON(w, http.StatusOK, json.SuccessResponse[any]{Data: resp.AccessToken})

}
