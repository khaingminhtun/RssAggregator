package posts

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
