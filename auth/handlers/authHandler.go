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

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dtos.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if req.Email == "" || req.Password == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "email and password required")
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
