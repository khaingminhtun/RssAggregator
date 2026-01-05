package feeds

import (
	"net/http"

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
