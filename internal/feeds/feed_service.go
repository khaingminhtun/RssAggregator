package feeds

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/khaingminhtun/rssagg/internal/adapters/database/repo"
	"github.com/khaingminhtun/rssagg/internal/pkg/errorHandle"
	"github.com/khaingminhtun/rssagg/internal/pkg/log"
	"github.com/khaingminhtun/rssagg/internal/pkg/utils"
	"github.com/khaingminhtun/rssagg/internal/posts"
)

type FeedService interface {
	CreateFeed(ctx context.Context, siteURL string) (*FeedResponse, error)
	// FetchFeedPosts(ctx context.Context, feedID int32) ([]PostResponse, error)

	// Read
	GetFeedByID(ctx context.Context, id int32) (*FeedResponse, error)
	GetAllFeeds(ctx context.Context) ([]FeedResponse, error)
	GetFeedsByUserID(ctx context.Context, userID int32) ([]FeedResponse, error)
}

type service struct {
	feedRepo FeedRepository
	postRepo posts.PostRepository
	fetcher  *FetcherService
}

func NewFeedService(feedRepo FeedRepository, postRepo posts.PostRepository, fetcher *FetcherService) FeedService {
	return &service{
		feedRepo: feedRepo,
		postRepo: postRepo,
		fetcher:  fetcher,
	}
}

// create feed from site url
func (s *service) CreateFeed(ctx context.Context, siteURL string) (*FeedResponse, error) {
	// 1. Validate input
	if siteURL == "" {
		return nil, errorHandle.BadRequest("site_url is required" + siteURL)
	}

	// 2. Validate URL
	u, err := utils.ValidateURL(siteURL)
	if err != nil {
		return nil, errorHandle.BadRequest("invalid site_url" + siteURL)
	}
	log.Info("validated url", "url", u.String())

	// 3. Discover RSS feed
	feedURL, err := s.fetcher.DiscoverFeed(u.String())
	if err != nil {
		return nil, err
	}
	log.Info("discovered feed url", "feed_url", feedURL)

	// 3a. Fetch and parse RSS feed to get title/description
	parsedFeed, err := s.fetcher.FetchAndParse(feedURL)
	if err != nil {
		return nil, err
	}

	// 4. Save feed in DB
	feed, err := s.feedRepo.CreateFeed(ctx, repo.CreateFeedParams{
		FeedUrl: feedURL,
		Title:   parsedFeed.Title,
		Description: pgtype.Text{
			String: parsedFeed.Description,
			Valid:  parsedFeed.Description != "",
		},
		WebsiteUrl: siteURL,
	})
	if err != nil {
		return nil, err
	}

	log.Info("feed created successfully", "feed_id", feed)

	return &FeedResponse{
		Title:       feed.Title,
		Description: feed.Description.String,
		WebsiteUrl:  feed.WebsiteUrl,
		FeedURL:     feed.FeedUrl,
	}, nil
}

// fetch feeds
// func (s *service) FetchFeedPosts(ctx context.Context, feedID int32) ([]PostResponse, error) {

// 	feedURL, err := s.feedRepo.GetFeedURLByID(ctx, feedID)
// 	if err != nil {
// 		return nil, errorHandle.NotFound("feedURL " + feedURL + " is not found")
// 	}

// 	parsedFeed, err := s.fetcher.FetchAndParse(feedURL)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var responses []PostResponse

// 	//  Iterate items
// 	for _, item := range parsedFeed.Items {
// 		if item.Link == "" {
// 			continue
// 		}

// 		// Published date
// 		var publishedAt time.Time
// 		if item.PublishedParsed != nil {
// 			publishedAt = *item.PublishedParsed
// 		} else {
// 			publishedAt = time.Now()
// 		}

// 		// GUID fallback
// 		guid := item.GUID
// 		if guid == "" {
// 			guid = item.Link
// 		}

// 		post, err := s.postRepo.CreatePost(ctx, repo.CreatePostParams{
// 			FeedID: feedID,
// 			Title:  item.Title,
// 			Url:    item.Link,
// 			Description: pgtype.Text{
// 				String: item.Description,
// 				Valid:  item.Description != "",
// 			},
// 			PublishedAt: pgtype.Timestamptz{
// 				Time:  publishedAt,
// 				Valid: true,
// 			},
// 			Guid: pgtype.Text{
// 				String: guid,
// 				Valid:  true,
// 			},
// 		})

// 		if err != nil {
// 			return nil, err
// 		}

// 		responses = append(responses, PostResponse{
// 			ID:          post.ID,
// 			FeedID:      post.FeedID,
// 			Title:       post.Title,
// 			URL:         post.Url,
// 			Description: post.Description.String,
// 			PublishedAt: post.PublishedAt.Time,
// 			Guid:        post.Guid.String,
// 			CreatedAt:   post.CreatedAt.Time,
// 		})
// 	}

// 	return responses, nil
// }

// get feed by id
func (s *service) GetFeedByID(ctx context.Context, id int32) (*FeedResponse, error) {
	feed, err := s.feedRepo.GetFeedByID(ctx, id)
	if err != nil {
		log.Error("failed to get feed by id", "error", err)
		return nil, errorHandle.NotFound("feed not found")
	}

	return mapFeedToResponse(feed), nil
}

// get all feeds
func (s *service) GetAllFeeds(ctx context.Context) ([]FeedResponse, error) {

	feeds, err := s.feedRepo.GetAllFeeds(ctx)
	if err != nil {
		log.Error("failed to get all feeds", "error", err)
		return nil, err
	}

	return mapFeedsToResponses(feeds), nil
}

// get feed by userid
func (s *service) GetFeedsByUserID(ctx context.Context, userID int32) ([]FeedResponse, error) {
	feeds, err := s.feedRepo.GetFeedsByUserID(ctx, userID)
	if err != nil {
		log.Error("failed to get feeds by user id", "userID", userID, "error", err)
		return nil, err
	}

	return mapFeedsToResponses(feeds), nil
}

// ----------------- Mapping helpers ----------------
func mapFeedToResponse(feed repo.Feed) *FeedResponse {
	return &FeedResponse{
		ID:            feed.ID,
		Title:         feed.Title,
		WebsiteUrl:    feed.WebsiteUrl,
		FeedURL:       feed.FeedUrl,
		CreatedAt:     feed.CreatedAt.Time,
		LastFetchedAt: feed.LastFetchedAt.Time,
	}
}

func mapFeedsToResponses(feeds []repo.Feed) []FeedResponse {
	result := make([]FeedResponse, 0, len(feeds))
	for _, feed := range feeds {
		result = append(result, *mapFeedToResponse(feed))
	}
	return result
}
