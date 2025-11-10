package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/khaingminhtun/rssagg/api/users/serializer"
	"github.com/khaingminhtun/rssagg/api/users/services"
	"github.com/khaingminhtun/rssagg/utils"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(s *services.UserService) *UserHandler {
	return &UserHandler{service: s}
}

type signupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req signupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if req.Email == "" || req.Password == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "email and password required")
	}

	user, err := h.service.RegisterUser(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, serializer.SerializeUser(user))

}

// GetUserById handler
func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	user, err := h.service.GetUserById(r.Context(), utils.StringToUUID(idParam))
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, err.Error())
	}

	utils.RespondWithJSON(w, http.StatusOK, serializer.SerializeUser(user))
}
