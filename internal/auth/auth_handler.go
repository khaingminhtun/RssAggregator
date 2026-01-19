package auth

import (
	"net/http"

	"github.com/khaingminhtun/rssagg/internal/pkg/errorHandle"
	"github.com/khaingminhtun/rssagg/internal/pkg/json"
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

	json.RespondJSON(w, http.StatusCreated, "user registered successfulyy", nil)

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

	json.RespondJSON(w, http.StatusOK, "user log in successfully", resp.AccessToken)

}

// refresh token
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest

	if !json.DecodeJSON(w, r, &req) {
		return
	}

	resp, err := h.AuthService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		errorHandle.RespondHTTPError(w, err)
		return
	}

	json.RespondJSON(w, http.StatusOK, "refresh token successfully", resp.AccessToken)
}

// logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear the refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,                // set to true in production
		SameSite: http.SameSiteLaxMode, // keep same as login
		MaxAge:   -1,                   // negative MaxAge deletes the cookie
	})

	json.RespondJSON(w, http.StatusOK, "user logged out successfully", nil)
}

