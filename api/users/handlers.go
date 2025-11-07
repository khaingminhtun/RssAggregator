package users

import (
	"encoding/json"
	"net/http"

	"github.com/khaingminhtun/rssagg/utils"
)

type UserHandler struct {
	service *UserService
}

func NewUserHandler(s *UserService) *UserHandler {
	return &UserHandler{service: s}
}

// POST /users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.service.RegisterUser(r.Context(), req.Name, req.Email)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Cannot create user")
		return
	}
	// //use serializer before sending
	resp :=
		SerializeUser(user)

	utils.RespondWithJSON(w, http.StatusCreated, resp)
}
