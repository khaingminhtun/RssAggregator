package feeds

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/khaingminhtun/rssagg/internal/errorHandle"
	"github.com/khaingminhtun/rssagg/internal/json"
)

type FeedHandler struct {
	FeedService FeedService
}

func NewFeedHandler(feedService FeedService) *FeedHandler {
	return &FeedHandler{
		FeedService: feedService,
	}
}

func (h *FeedHandler) CreateFeed(w http.ResponseWriter, r *http.Request) {
	var req RequestSiteURL

	if !json.DecodeJSON(w, r, &req) {
		return
	}

	feed, err := h.FeedService.CreateFeed(r.Context(), req.SiteURL)
	if err != nil {
		errorHandle.RespondHTTPError(w, err)
		return
	}

	json.RespondJSON(w, http.StatusCreated, "feed create successfully", feed)
}

// GET /feeds/{feedID}/posts
func (h *FeedHandler) GetFeedPosts(w http.ResponseWriter, r *http.Request) {
	//1. Parse feedID from URL
	feedIDStr := chi.URLParam(r, "feedID")
	if feedIDStr == "" {
		errorHandle.RespondHTTPError(w, errorHandle.ErrInvalidInput)
		return
	}

	feedID, err := strconv.ParseInt(feedIDStr, 10, 32)
	if err != nil {
		errorHandle.RespondHTTPError(w, errorHandle.ErrInvalidInput)
		return
	}

	//2. call service
	posts, err := h.FeedService.FetchFeedPosts(r.Context(), int32(feedID))
	if err != nil {
		errorHandle.RespondHTTPError(w, err)
		return
	}

	json.RespondJSON(w, http.StatusOK, "feed posts fetched successfully", posts)
}
