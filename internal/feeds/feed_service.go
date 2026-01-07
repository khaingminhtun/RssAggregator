package feeds

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/khaingminhtun/rssagg/internal/adapters/database/repo"
	"github.com/khaingminhtun/rssagg/internal/errorHandle"
	"github.com/khaingminhtun/rssagg/internal/log"
	"github.com/khaingminhtun/rssagg/internal/utils"
)

type FeedService interface {
	CreateFeed(ctx context.Context, siteURL string) (*FeedResponse, error)
	FetchFeedPosts(ctx context.Context, feedID int32) ([]PostResponse, error)
}

type service struct {
	repo    repo.Querier
	fetcher *FetcherService
}

func NewFeedService(repo repo.Querier, fetcher *FetcherService) FeedService {
	return &service{
		repo:    repo,
		fetcher: fetcher,
	}
}

// create feed from site url
func (s *service) CreateFeed(ctx context.Context, siteURL string) (*FeedResponse, error) {
	// 1. Validate input
	if siteURL == "" {
		return nil, errorHandle.ErrInvalidInput
	}

	// 2. Validate URL
	u, err := utils.ValidateURL(siteURL)

	if err != nil {
		return nil, errorHandle.ErrInvalidCredentials
	}
	log.Info("validated url", "url", u.String())

	// 3. Discover RSS feed
	feedURL, err := s.fetcher.DiscoverFeed(u.String())
	if err != nil {
		return nil, err
	}
	log.Info("discovered feed url", "feed_url", feedURL)

	// 4. Save feed in DB
	feed, err := s.repo.CreateFeed(ctx, repo.CreateFeedParams{
		Url:         feedURL,
		Title:       "ok",
		Description: utils.Text("okay"),
		SiteUrl:     siteURL,
	})
	if err != nil {
		if errorHandle.IsUniqueViolation(err) {
			return nil, errorHandle.ErrUserAlreadyExists
		}
		return nil, errorHandle.ErrInternal
	}

	log.Info("feed created successfully", "feed_id", feed)

	return &FeedResponse{
		Title:   feed.Title,
		SiteURL: feed.SiteUrl,
		FeedURL: feed.Url,
	}, nil
}

// fetch feeds
func (s *service) FetchFeedPosts(ctx context.Context, feedID int32) ([]PostResponse, error) {
	// 1. Get feed URL
	feedURL, err := s.repo.GetFeedURLByID(ctx, feedID)
	if err != nil {
		return nil, errorHandle.ErrUserNotFound
	}

	// 2. Fetch & parse RSS
	parsedFeed, err := s.fetcher.FetchAndParse(feedURL)
	if err != nil {
		return nil, errorHandle.ErrInternal
	}

	var responses []PostResponse

	// 3. Iterate items
	for _, item := range parsedFeed.Items {
		if item.Link == "" {
			continue
		}

		// Published date
		var publishedAt time.Time
		if item.PublishedParsed != nil {
			publishedAt = *item.PublishedParsed
		} else {
			publishedAt = time.Now()
		}

		// GUID fallback
		guid := item.GUID
		if guid == "" {
			guid = item.Link
		}

		post, err := s.repo.CreatePost(ctx, repo.CreatePostParams{
			FeedID: feedID,
			Title:  item.Title,
			Url:    item.Link,
			Description: pgtype.Text{
				String: item.Description,
				Valid:  item.Description != "",
			},
			PublishedAt: pgtype.Timestamptz{
				Time:  publishedAt,
				Valid: true,
			},
			Guid: pgtype.Text{
				String: guid,
				Valid:  true,
			},
		})

		if err != nil {
			if errorHandle.IsUniqueViolation(err) {
				continue
			}
			log.Error("failed to insert post", "feed_id", feedID, "err", err)
			continue
		}

		responses = append(responses, PostResponse{
			ID:          post.ID,
			FeedID:      post.FeedID,
			Title:       post.Title,
			URL:         post.Url,
			Description: post.Description.String,
			PublishedAt: post.PublishedAt.Time,
			Guid:        post.Guid.String,
			CreatedAt:   post.CreatedAt.Time,
		})
	}

	return responses, nil
}
