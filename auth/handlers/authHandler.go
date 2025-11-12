package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	apiResponse "github.com/khaingminhtun/rssagg/api/dtos"
	userDtos "github.com/khaingminhtun/rssagg/api/users/dtos"
	"github.com/khaingminhtun/rssagg/auth/dtos"
	"github.com/khaingminhtun/rssagg/auth/services"
	"github.com/khaingminhtun/rssagg/utils"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(s *services.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

// register
func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dtos.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "name and email and password required")
	}

	user, err := a.service.RegisterUser(r.Context(), req)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, apiResponse.APIResponse{
		Success: true,
		Data:    userDtos.SerializeUser(user),
		Message: fmt.Sprintf("User %s registered successfully", user.Name),
	})

}

// Auth
func (a *AuthHandler) Authenicate(w http.ResponseWriter, r *http.Request) {
	var req dtos.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Authenticate and generate tokens
	resp, err := a.service.AuthenticateUser(r.Context(), req)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "user not authorized")
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

	utils.RespondWithJSON(w, http.StatusOK, apiResponse.APIResponse{
		Success: true,
		Data:    resp.AccessToken,

		Message: fmt.Sprintf("User with email %s login successfully", req.Email),
	})

}

// refreshtoken
func (a *AuthHandler) GetRefreshToken(w http.ResponseWriter, r *http.Request) {
	// 1️⃣ Get refresh token from cookie
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "refresh token not found")
		return
	}
	refreshToken := cookie.Value

	// 2️⃣ Call service
	resp, err := a.service.RefreshToken(r.Context(), refreshToken)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// set new refresh token in cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    resp.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // set true in production
		SameSite: http.SameSiteLaxMode,
		MaxAge:   7 * 24 * 60 * 60,
	})

	utils.RespondWithJSON(w, http.StatusOK, apiResponse.APIResponse{
		Success: true,
		Data:    resp,
		Message: "get refreshtoken and accesstoken return  successfully",
	})
}

// logout
func (a *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// 1️⃣ Get refresh token from cookie
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "no refresh token not found")
		return
	}

	// 2️⃣ Call service to delete refresh token from Redis
	err = a.service.LogoutUser(r.Context(), cookie.Value)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to logout")
		return
	}

	// 3️⃣ Delete cookie in browser
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false, // true in prod
	})

	utils.RespondWithJSON(w, http.StatusOK, "logged out successfully")
}
