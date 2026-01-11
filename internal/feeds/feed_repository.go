package feeds

import (
	"context"

	"github.com/khaingminhtun/rssagg/internal/adapters/database/repo"
)

type FeedRepository interface {
	// Create
	CreateFeed(ctx context.Context, arg repo.CreateFeedParams) (repo.Feed, error)

	// Read
	GetFeedByID(ctx context.Context, id int32) (repo.Feed, error)
	GetAllFeeds(ctx context.Context) ([]repo.Feed, error)
	GetFeedsByUserID(ctx context.Context, userID int32) ([]repo.Feed, error)
	GetFeedURLByID(ctx context.Context, id int32) (string, error)

	// Update
	UpdateFeed(ctx context.Context, arg repo.UpdateFeedParams) (repo.Feed, error)

	// Delete
	DeleteFeedByID(ctx context.Context, id int32) error
	DeleteFeedIfUnused(ctx context.Context, id int32) error
	DeleteAllFeedsByUserID(ctx context.Context, userID int32) error
}
