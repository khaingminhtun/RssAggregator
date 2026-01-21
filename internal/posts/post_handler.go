package posts

import (
	"net/http"

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
