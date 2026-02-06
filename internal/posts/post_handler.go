package posts

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/khaingminhtun/rssagg/internal/auth"
	"github.com/khaingminhtun/rssagg/internal/pkg/errorHandle"
	"github.com/khaingminhtun/rssagg/internal/pkg/json"
)

type PostHandler struct {
	PostService PostService
}

func NewPostHandler(postService PostService) *PostHandler {
	return &PostHandler{
		PostService: postService,
	}
}

// GET /posts - get all posts
func (h *PostHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.PostService.GetAllPosts(r.Context())
	if err != nil {
		errorHandle.RespondHTTPsError(w, r, err)
		return
	}

	json.RespondJSON(
		w,
		http.StatusOK,
		"all posts fetched successfully",
		posts,
	)
}

// GET /posts/{id}
func (h *PostHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")
	if postID == "" {
		errorHandle.RespondHTTPsError(
			w,
			r,
			errorHandle.BadRequest("post id is required"),
		)
		return
	}

	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		errorHandle.RespondHTTPsError(
			w,
			r,
			errorHandle.BadRequest("post id is not a valid integer"),
		)
		return
	}

	post, err := h.PostService.GetPostByID(r.Context(), int32(postIDInt))
	if err != nil {
		errorHandle.RespondHTTPsError(w, r, err)
		return
	}
	json.RespondJSON(
		w,
		http.StatusOK,
		"post fetched successfully",
		post,
	)
}

// GET /feeds/{feedID}/posts
func (h *PostHandler) GetPostsByFeedID(w http.ResponseWriter, r *http.Request) {
	feedID := chi.URLParam(r, "feedID")
	if feedID == "" {
		errorHandle.RespondHTTPsError(
			w,
			r,
			errorHandle.BadRequest("feedID path parameter is required"),
		)
		return
	}

	feedIDInt, err := strconv.Atoi(feedID)
	if err != nil {
		errorHandle.RespondHTTPsError(
			w,
			r,
			errorHandle.BadRequest("feedID is not a valid integer"),
		)
		return
	}

	posts, err := h.PostService.GetPostsByFeedID(r.Context(), int32(feedIDInt))
	if err != nil {
		errorHandle.RespondHTTPsError(w, r, err)
		return
	}
	json.RespondJSON(
		w,
		http.StatusOK,
		"posts fetched successfully for feed",
		posts,
	)
}

// GET /posts/latest?limit=10 - get latest posts with limit
func (h *PostHandler) GetLatestPosts(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit, _ := strconv.Atoi(limitStr)

	posts, err := h.PostService.GetLatestPosts(
		r.Context(),
		int32(limit),
	)
	if err != nil {
		errorHandle.RespondHTTPsError(w, r, err)
		return
	}

	json.RespondJSON(w, http.StatusOK, "latest posts fetched successfully", posts)
}

// GET /posts/search?query=example&limit=10&offset=0 - search posts by query with pagination
func (h *PostHandler) SearchPosts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	posts, err := h.PostService.SearchPosts(
		r.Context(),
		q,
		int32(page),
		int32(limit),
	)
	if err != nil {
		errorHandle.RespondHTTPError(w, err)
		return
	}

	json.RespondJSON(w, http.StatusOK, "posts searched successfully", posts)
}

// GET /posts/timeline?page=1&limit=20
// Handler to get all posts from feeds the current user is subscribed to
func (h *PostHandler) GetTimeline(w http.ResponseWriter, r *http.Request) {
	// Assuming userID is stored in context after auth middleware
	userIDVal := r.Context().Value("userID")
	userID, ok := userIDVal.(int32)
	if !ok || userID <= 0 {
		errorHandle.RespondHTTPError(w, fmt.Errorf("unauthorized or invalid user"))
		return
	}

	// Pagination query params
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	// Call service
	posts, err := h.PostService.GetPostsForUser(r.Context(), userID, int32(page), int32(limit))
	if err != nil {
		errorHandle.RespondHTTPError(w, err)
		return
	}

	// Respond JSON
	json.RespondJSON(w, http.StatusOK, "timeline fetched successfully", posts)
}

// PATCH /posts/{id}/read - mark post as read/unread for user
func (h *PostHandler) MarkPostRead(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		json.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	postID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		json.RespondError(w, http.StatusBadRequest, "invalid post id")
		return
	}

	var req MarkReadRequest
	// 1. Decode request body
	if !json.DecodeJSON(w, r, &req) {
		return
	}

	if err := h.PostService.MarkPostRead(
		r.Context(),
		userID,
		int32(postID),
		req.IsRead,
	); err != nil {
		json.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.RespondJSON(w, http.StatusOK, "post read state updated",
		nil)
}

// PATCH /posts/{id}/favorite - mark post as favourite/unfavourite for user
func (h *PostHandler) MarkPostFavourite(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		json.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	postID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		json.RespondError(w, http.StatusBadRequest, "invalid post id")
		return
	}

	var req MarkFavouriteRequest
	// 1. Decode request body
	if !json.DecodeJSON(w, r, &req) {
		return
	}

	if err := h.PostService.MarkPostFavourite(
		r.Context(),
		userID,
		int32(postID),
		req.IsFavourite,
	); err != nil {
		json.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.RespondJSON(w, http.StatusOK,
		"post favourite state updated", nil)
}
