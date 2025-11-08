package users

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.service.RegisterUser(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Cannot create user")
		return
	}
	// //use serializer before sending
	resp :=
		SerializeUser(user)

	utils.RespondWithJSON(w, http.StatusCreated, resp)
}

// Get Usersby Id
func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.service.GetUserById(r.Context(), int32(id))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "User not found")
		return
	}

	resp := SerializeUser(user)
	utils.RespondWithJSON(w, http.StatusOK, resp)
}
