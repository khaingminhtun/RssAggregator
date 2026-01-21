package mechanism

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/khaingminhtun/rssagg/internal/adapters/database/repo"
	"github.com/khaingminhtun/rssagg/internal/feeds"
	"github.com/khaingminhtun/rssagg/internal/posts"
)

// ----------------------------
// FeedProcessor interface
// ----------------------------
type FeedProcessor interface {
	ProcessFeed(ctx context.Context, feed repo.Feed) error
}

// ----------------------------
// Service implementation
// ----------------------------
type FeedProcessorService struct {
	fetcher  *feeds.FetcherService
	postRepo posts.PostRepository
	feedRepo feeds.FeedRepository
}

func NewFeedProcessorService(fetcher *feeds.FetcherService, postRepo posts.PostRepository, feedRepo feeds.FeedRepository) *FeedProcessorService {
	return &FeedProcessorService{
		fetcher:  fetcher,
		postRepo: postRepo,
		feedRepo: feedRepo,
	}
}

// ----------------------------
// Process a feed
// ----------------------------
func (s *FeedProcessorService) ProcessFeed(ctx context.Context, feed repo.Feed) error {
	if feed.FeedUrl == "" {
    log.Printf("[Processor] Skipping feed %d: empty FeedUrl", feed.ID)
    return nil
}

	parsedFeed, err := s.fetcher.FetchAndParse(feed.FeedUrl)
	if err != nil {
		return err
	}

	for _, item := range parsedFeed.Items {
		if item.Link == "" {
			continue
		}

		publishedAt := time.Now()
		if item.PublishedParsed != nil {
			publishedAt = *item.PublishedParsed
		}

		guid := item.GUID
		if guid == "" {
			guid = item.Link
		}

		_, err = s.postRepo.CreatePost(ctx, repo.CreatePostParams{
			FeedID: feed.ID,
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
			log.Printf("Failed to create post %s: %v", item.Title, err)
			return err
		}
	}

	return nil
}
