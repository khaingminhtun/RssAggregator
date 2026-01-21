package feeds

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/khaingminhtun/rssagg/internal/pkg/errorHandle"
	"github.com/khaingminhtun/rssagg/internal/pkg/json"
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

// // GET /feeds/{feedID}/posts
// func (h *FeedHandler) GetFeedPosts(w http.ResponseWriter, r *http.Request) {
// 	//1. Parse feedID from URL
// 	feedIDStr := chi.URLParam(r, "feedID")
// 	if feedIDStr == "" {
// 		errorHandle.RespondHTTPsError(w, r, errorHandle.BadRequest("feedID is required"))
// 		return
// 	}

// 	feedID, err := strconv.ParseInt(feedIDStr, 10, 32)
// 	if err != nil {
// 		errorHandle.RespondHTTPsError(w, r, errorHandle.BadRequest("feedID must be a number"))
// 		return
// 	}

// 	//2. call service
// 	posts, err := h.FeedService.FetchFeedPosts(r.Context(), int32(feedID))
// 	if err != nil {
// 		errorHandle.RespondHTTPsError(w, r, err)
// 		return
// 	}

// 	json.RespondJSON(w, http.StatusOK, "feed posts fetched successfully", posts)
// }

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

// Update /feeds/{id} - update feed
// update feed
func (h *FeedHandler) UpdateFeed(w http.ResponseWriter, r *http.Request) {
	var req RequestSiteURL

	// 1. Decode request body
	if !json.DecodeJSON(w, r, &req) {
		return
	}

	// 2. Get feed ID from URL
	idParam := chi.URLParam(r, "id")
	feedID, err := strconv.Atoi(idParam)
	if err != nil || feedID <= 0 {
		errorHandle.RespondHTTPError(w, errorHandle.BadRequest("invalid feed id"))
		return
	}

	// 3. Call service
	feed, err := h.FeedService.UpdateFeed(
		r.Context(),
		int32(feedID),
		req.SiteURL,
	)
	if err != nil {
		errorHandle.RespondHTTPError(w, err)
		return
	}

	// 4. Respond
	json.RespondJSON(w, http.StatusOK, "feed updated successfully", feed)
}

// DELETE /feeds/{id}
func (h *FeedHandler) DeleteFeedByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	feedID, err := strconv.Atoi(idParam)
	if err != nil || feedID <= 0 {
		errorHandle.RespondHTTPError(w, errorHandle.BadRequest("invalid feed id"))
		return
	}

	err = h.FeedService.DeleteFeedByID(r.Context(), int32(feedID))
	if err != nil {
		errorHandle.RespondHTTPError(w, err)
		return
	}

	json.RespondJSON(w, http.StatusOK, "feed deleted successfully", nil)
}

// DELETE /feeds/{id}/unused
func (h *FeedHandler) DeleteFeedIfUnused(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	feedID, err := strconv.Atoi(idParam)
	if err != nil || feedID <= 0 {
		errorHandle.RespondHTTPError(w, errorHandle.BadRequest("invalid feed id"))
		return
	}

	err = h.FeedService.DeleteFeedIfUnused(r.Context(), int32(feedID))
	if err != nil {
		errorHandle.RespondHTTPError(w, err)
		return
	}

	json.RespondJSON(w, http.StatusOK, "unused feed deleted successfully", nil)
}

// DELETE /users/{user_id}/feeds
func (h *FeedHandler) DeleteAllFeedsByUserID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "user_id")
	userID, err := strconv.Atoi(idParam)
	if err != nil || userID <= 0 {
		errorHandle.RespondHTTPError(w, errorHandle.BadRequest("invalid user id"))
		return
	}

	err = h.FeedService.DeleteAllFeedsByUserID(r.Context(), int32(userID))
	if err != nil {
		errorHandle.RespondHTTPError(w, err)
		return
	}

	json.RespondJSON(w, http.StatusOK, "all feeds deleted successfully", nil)
}
