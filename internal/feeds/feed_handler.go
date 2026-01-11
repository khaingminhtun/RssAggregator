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

// create feed
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
		errorHandle.RespondHTTPsError(w, r, errorHandle.BadRequest("feedID is required"))
		return
	}

	feedID, err := strconv.ParseInt(feedIDStr, 10, 32)
	if err != nil {
		errorHandle.RespondHTTPsError(w, r, errorHandle.BadRequest("feedID must be a number"))
		return
	}

	//2. call service
	posts, err := h.FeedService.FetchFeedPosts(r.Context(), int32(feedID))
	if err != nil {
		errorHandle.RespondHTTPsError(w, r, err)
		return
	}

	json.RespondJSON(w, http.StatusOK, "feed posts fetched successfully", posts)
}

// GET /feeds/{feedID} - get feed by ID
func (h *FeedHandler) GetFeedByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		errorHandle.RespondHTTPsError(w, r, errorHandle.BadRequest("feed ID is required"))
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		errorHandle.RespondHTTPsError(w, r, errorHandle.BadRequest("feed ID must be a number"))
		return
	}

	feed, err := h.FeedService.GetFeedByID(r.Context(), int32(id))
	if err != nil {
		errorHandle.RespondHTTPsError(w, r, err)
		return
	}

	json.RespondJSON(w, http.StatusOK, "feed fetched successfully", feed)
}

// GET /feeds - get all feeds
func (h *FeedHandler) GetAllFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := h.FeedService.GetAllFeeds(r.Context())
	if err != nil {
		errorHandle.RespondHTTPsError(w, r, err)
		return
	}

	json.RespondJSON(w, http.StatusOK, "all feeds fetched successfully", feeds)
}

// GET /users/{userID}/feeds - get feeds by user ID
func (h *FeedHandler) GetFeedsByUserID(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	if userIDStr == "" {
		errorHandle.RespondHTTPsError(w, r, errorHandle.BadRequest("user ID is required"))
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		errorHandle.RespondHTTPsError(w, r, errorHandle.BadRequest("user ID must be a number"))
		return
	}

	feeds, err := h.FeedService.GetFeedsByUserID(r.Context(), int32(userID))
	if err != nil {
		errorHandle.RespondHTTPsError(w, r, err)
		return
	}

	json.RespondJSON(w, http.StatusOK, "feeds fetched successfully for user", feeds)
}
