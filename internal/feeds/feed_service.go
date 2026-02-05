package feeds

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/khaingminhtun/rssagg/internal/adapters/cache"
	"github.com/khaingminhtun/rssagg/internal/adapters/database/repo"
	"github.com/khaingminhtun/rssagg/internal/pkg/errorHandle"
	"github.com/khaingminhtun/rssagg/internal/pkg/log"
	"github.com/khaingminhtun/rssagg/internal/pkg/utils"
	"github.com/khaingminhtun/rssagg/internal/posts"
)

type FeedService interface {
	CreateFeed(ctx context.Context, siteURL string) (*FeedResponse, error)

	SubscribeUserToFeed(ctx context.Context, userID int32, feedID int32) error

	// Read
	GetFeedByID(ctx context.Context, id int32) (*FeedResponse, error)
	GetAllFeeds(ctx context.Context) ([]FeedResponse, error)
	GetFeedsByUserID(ctx context.Context, userID int32) ([]FeedResponse, error)
	

	// Update
	UpdateFeed(ctx context.Context, feedID int32, siteUrl string) (*FeedResponse, error)

	// Delete
	DeleteFeedByID(ctx context.Context, feedID int32) error
	DeleteFeedIfUnused(ctx context.Context, feedID int32) error
	DeleteAllFeedsByUserID(ctx context.Context, userID int32) error
}

type service struct {
	feedRepo FeedRepository
	postRepo posts.PostRepository
	fetcher  *FetcherService
	rssCache *cache.RSSURLCache
}

func NewFeedService(feedRepo FeedRepository, postRepo posts.PostRepository, fetcher *FetcherService, rssCache *cache.RSSURLCache) FeedService {
	return &service{
		feedRepo: feedRepo,
		postRepo: postRepo,
		fetcher:  fetcher,
		rssCache: rssCache,
	}
}

// create feed from site url
// create feed from site url
func (s *service) CreateFeed(ctx context.Context, siteURL string) (*FeedResponse, error) {
	// 1. Validate input
	if siteURL == "" {
		return nil, errorHandle.BadRequest("site_url is required")
	}

	// 2. Validate & normalize URL
	u, err := utils.ValidateURL(siteURL)
	if err != nil {
		return nil, errorHandle.BadRequest("invalid site_url")
	}
	normalizedURL := u.String()
	log.Info("validated url", "url", normalizedURL)

	var feedURL string

	// 3. Try Redis cache
	if cachedURL, ok, err := s.rssCache.Get(ctx, normalizedURL); err != nil {
		log.Info("rss cache unavailable", "error", err)
	} else if ok {
		feedURL = cachedURL
		log.Info("rss url cache hit", "feed_url", feedURL)
	}

	// 4. Cache miss → discover feed
	if feedURL == "" {
		discoveredURL, err := s.fetcher.DiscoverFeed(ctx, normalizedURL)
		if err != nil {
			log.Error("feed discovery failed", "url", normalizedURL, "error", err)
			return nil, errorHandle.FeedDiscoveryFailed(
				"unable to discover RSS feed from site",
			)
		}
		feedURL = discoveredURL
		log.Info("rss feed discovered", "feed_url", feedURL)
	}

	// 5. Fetch & parse (authoritative validation step)
	parsedFeed, err := s.fetcher.FetchAndParse(ctx, feedURL)
	if err != nil {
		log.Error("feed parsing failed", "feed_url", feedURL, "error", err)

		// Cache poisoning protection
		if feedURL != "" {
			_ = s.rssCache.Delete(ctx, normalizedURL)
		}

		return nil, errorHandle.FeedParseFailed(
			"failed to fetch or parse RSS feed",
		)
	}

	// 6. Cache ONLY after full success
	if err := s.rssCache.Set(ctx, normalizedURL, feedURL); err != nil {
		log.Info("failed to cache rss url", "error", err)
	}

	// 7. Persist in DB
	feed, err := s.feedRepo.CreateFeed(ctx, repo.CreateFeedParams{
		FeedUrl:    feedURL,
		Title:      parsedFeed.Title,
		WebsiteUrl: normalizedURL,
		Description: utils.NullString(parsedFeed.Description),
	})
	if err != nil {
		log.Error("feed creation failed in DB", "feed_url", feedURL, "error", err)
		return nil, errorHandle.FeedAlreadyExists(
			"feed with this URL already exists",
		)
	}

	// userID, ok := auth.UserIDFromContext(ctx)
	// if ok && userID > 0 {
	// 	err = s.SubscribeUserToFeed(ctx, userID, feed.ID)
	// 	if err != nil {
	// 		log.Error("failed to subscribe user to feed after creation", "user_id", userID, "feed_id", feed.ID, "error", err)
	// 		return nil, err
	// 	}
	// }

	const layout = "2006-01-02 15:04:05"

	return &FeedResponse{
		ID:            feed.ID,
		Title:         feed.Title,
		Description:   feed.Description.String,
		WebsiteUrl:    feed.WebsiteUrl,
		FeedURL:       feed.FeedUrl,
		CreatedAt:     feed.CreatedAt.Format(layout),
		LastFetchedAt: feed.LastFetchedAt.Time.Format(layout),
	}, nil
}

// subscribe user to feed
func (s *service) SubscribeUserToFeed(ctx context.Context, userID int32, feedID int32) error {

	// 1. Verify feed exists
	feed, err := s.feedRepo.GetFeedByID(ctx, feedID)
	if err != nil {
		log.Error("feed not found", "feed_id", feedID, "error", err)
		return errorHandle.NotFound("feed does not exist")
	}

	_, err = s.feedRepo.SubscribeUserToFeed(ctx, userID, feedID)
	if err != nil {
		log.Error("failed to subscribe user to feed", "user_id", userID, "feed_id", feedID, "error", err)
		return errorHandle.IsAlreadySubscribed("user is already subscribed to this feed")
	}

	log.Info("user subscribed to feed", "user_id", userID, "feed_id", feed.ID)
	return nil
}

// get feed by id
func (s *service) GetFeedByID(ctx context.Context, id int32) (*FeedResponse, error) {
	feed, err := s.feedRepo.GetFeedByID(ctx, id)
	if err != nil {
		log.Error("failed to get feed", "error", err)
		return nil, errorHandle.DatabaseError("unable to fetch feeds from DB")
	}
	return mapFeedToResponse(feed), nil
}

// get all feeds
func (s *service) GetAllFeeds(ctx context.Context) ([]FeedResponse, error) {

	feeds, err := s.feedRepo.GetAllFeeds(ctx)
	if err != nil {
		log.Error("failed to get all feeds", "error", err)
		return nil, errorHandle.DatabaseError("unable to fetch feeds from DB")
	}

	return mapFeedsToResponses(feeds), nil
}

// get feed by userid
func (s *service) GetFeedsByUserID(ctx context.Context, userID int32) ([]FeedResponse, error) {
	if userID <= 0 {
		return nil, errorHandle.BadRequest("invalid user_id")
	}

	feeds, err := s.feedRepo.GetFeedsByUserID(ctx, userID)
	if err != nil {
		log.Error("failed to get feeds by user id", "userID", userID, "error", err)
		return nil, errorHandle.DatabaseError("unable to fetch user's feeds")
	}

	return mapFeedsToResponses(feeds), nil
}

// update feed
func (s *service) UpdateFeed(
	ctx context.Context,
	feedID int32,
	siteURL string,
) (*FeedResponse, error) {

	// 1. Validate input
	if feedID <= 0 {
		return nil, errorHandle.BadRequest("invalid feed_id")
	}

	if siteURL == "" {
		return nil, errorHandle.BadRequest("site_url is required")
	}

	// 2. Validate URL
	u, err := utils.ValidateURL(siteURL)
	if err != nil {
		return nil, errorHandle.BadRequest("invalid site_url")
	}
	log.Info("validated url", "url", u.String())

	// 3. Discover RSS feed again (site may have changed)
	feedURL, err := s.fetcher.DiscoverFeed(ctx, u.String())
	if err != nil {
		log.Error("feed discovery failed", "url", u.String(), "error", err)
		return nil, errorHandle.FeedDiscoveryFailed("unable to discover RSS feed")
	}
	log.Info("discovered feed url", "feed_url", feedURL)

	// 4. Fetch & parse RSS feed
	parsedFeed, err := s.fetcher.FetchAndParse(ctx, feedURL)
	if err != nil {
		log.Error("feed parsing failed", "feed_url", feedURL, "error", err)
		return nil, errorHandle.FeedParseFailed("failed to parse RSS feed")
	}

	// 5. Update feed in DB
	updatedFeed, err := s.feedRepo.UpdateFeed(ctx, repo.UpdateFeedParams{
		ID:      feedID,
		Title:   parsedFeed.Title,
		FeedUrl: feedURL,
		Description: utils.NullString(parsedFeed.Description),
		WebsiteUrl: siteURL,
		LastFetchedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	})
	if err != nil {
		log.Error("feed update failed in DB", "feed_id", feedID, "error", err)
		return nil, errorHandle.NotFound("feed not found or update failed")
	}

	log.Info("feed updated successfully", "feed_id", updatedFeed.ID)

	// 6. Map response
	return &FeedResponse{
		Title:       updatedFeed.Title,
		Description: updatedFeed.Description.String,
		WebsiteUrl:  updatedFeed.WebsiteUrl,
		FeedURL:     updatedFeed.FeedUrl,
	}, nil
}

// delete
// ----------------- Delete Feed By ID -----------------
func (s *service) DeleteFeedByID(ctx context.Context, feedID int32) error {
	if feedID <= 0 {
		return errorHandle.BadRequest("invalid feed_id")
	}

	err := s.feedRepo.DeleteFeedByID(ctx, feedID)
	if err != nil {
		log.Error("failed to delete feed", "feed_id", feedID, "error", err)
		return errorHandle.DatabaseError("failed to delete feed")
	}

	log.Info("feed deleted", "feed_id", feedID)
	return nil
}

// ----------------- Delete Feed If Unused -----------------
func (s *service) DeleteFeedIfUnused(ctx context.Context, feedID int32) error {
	if feedID <= 0 {
		return errorHandle.BadRequest("invalid feed_id")
	}

	err := s.feedRepo.DeleteFeedIfUnused(ctx, feedID)
	if err != nil {
		log.Error("failed to delete unused feed", "feed_id", feedID, "error", err)
		return errorHandle.DatabaseError("failed to delete unused feed")
	}

	log.Info("unused feed deleted", "feed_id", feedID)
	return nil
}

// ----------------- Delete All Feeds By User -----------------
func (s *service) DeleteAllFeedsByUserID(ctx context.Context, userID int32) error {
	if userID <= 0 {
		return errorHandle.BadRequest("invalid user_id")
	}

	err := s.feedRepo.DeleteAllFeedsByUserID(ctx, userID)
	if err != nil {
		log.Error("failed to delete all feeds for user", "user_id", userID, "error", err)
		return errorHandle.DatabaseError("failed to delete user's feeds")
	}

	log.Info("all feeds deleted for user", "user_id", userID)
	return nil
}

// ----------------- Mapping helpers ----------------
func mapFeedToResponse(feed repo.Feed) *FeedResponse {

	const layout = "2006-01-02 15:04:05" // custom format
	return &FeedResponse{
		ID:            feed.ID,
		Title:         feed.Title,
		WebsiteUrl:    feed.WebsiteUrl,
		FeedURL:       feed.FeedUrl,
		CreatedAt:     feed.CreatedAt.Format(layout),
		LastFetchedAt: feed.LastFetchedAt.Time.Format(layout),
	}
}

func mapFeedsToResponses(feeds []repo.Feed) []FeedResponse {
	result := make([]FeedResponse, 0, len(feeds))
	for _, feed := range feeds {
		result = append(result, *mapFeedToResponse(feed))
	}
	return result
}
