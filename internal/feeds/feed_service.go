package feeds

import (
	"context"

	"github.com/khaingminhtun/rssagg/internal/adapters/database/repo"
	"github.com/khaingminhtun/rssagg/internal/errorHandle"
	"github.com/khaingminhtun/rssagg/internal/log"
	"github.com/khaingminhtun/rssagg/internal/utils"
)

type FeedService interface {
	CreateFeed(ctx context.Context, siteURL string) (*FeedResponse, error)
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

func (s *service) CreateFeed(ctx context.Context, siteURL string) (*FeedResponse, error) {
	// 1. Validate input
	if siteURL == "" {
		return nil, errorHandle.ErrInvalidInput
	}

	// 2. Validate URL
	u, err := utils.ValidateURL(siteURL)
	if err != nil {
		return nil, err
	}

	// 3. Discover RSS feed
	feedURL, err := s.fetcher.DiscoverFeed(u.String())
	if err != nil {
		return nil, err
	}

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
