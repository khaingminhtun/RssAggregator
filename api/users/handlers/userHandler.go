package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/khaingminhtun/rssagg/api/users/dtos"
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

// GetUserById handler
func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	user, err := h.service.GetUserById(r.Context(), utils.StringToUUID(idParam))
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, err.Error())
	}

	utils.RespondWithJSON(w, http.StatusOK, dtos.SerializeUser(user))
}
