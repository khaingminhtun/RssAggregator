package users

import (
	"net/http"
	"strconv"

	"github.com/khaingminhtun/rssagg/internal/auth"
	"github.com/khaingminhtun/rssagg/internal/pkg/json"
	"github.com/khaingminhtun/rssagg/internal/posts"
)

type UserHandler struct {
	PostService posts.PostService
}

func NewUserHandler(postService posts.PostService) *UserHandler {
	return &UserHandler{
		PostService: postService,
	}
}

// GET /users/me/favourites - get favourite posts for user
func (h *UserHandler) GetFavoritePosts(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		json.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	posts, err := h.PostService.GetFavouritePostsForUser(
		r.Context(),
		userID,
		int32(page),
		int32(limit),
	)
	if err != nil {
		json.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.RespondJSON(w, http.StatusOK, "favourite posts fetched successfully", posts)
}

// GET /users/me/posts - get all posts for user (read/unread)
func (h *UserHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		json.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	posts, err := h.PostService.GetPostsForUser(r.Context(), userID, int32(page), int32(limit))
	if err != nil {
		json.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.RespondJSON(w, http.StatusOK, "posts fetched successfully", posts)
}
