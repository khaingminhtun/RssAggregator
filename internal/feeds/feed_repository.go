package feeds

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/khaingminhtun/rssagg/internal/adapters/database/repo"
)

// FeedRepository is the domain-level interface your service will use.
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
	// Update only last fetched timestamp
	UpdateFeedFetchTime(ctx context.Context, id int32, lastFetchedAt pgtype.Timestamptz) (repo.Feed, error)

	// Delete
	DeleteFeedByID(ctx context.Context, id int32) error
	DeleteFeedIfUnused(ctx context.Context, id int32) error
	DeleteAllFeedsByUserID(ctx context.Context, userID int32) error
}

// --- Concrete SQLC-backed implementation --- //

type feedRepo struct {
	q *repo.Queries
}

// NewFeedRepository returns a FeedRepository implemented with SQLC queries
func NewFeedRepository(db repo.DBTX) FeedRepository {
	return &feedRepo{
		q: repo.New(db),
	}
}

func (r *feedRepo) CreateFeed(ctx context.Context, arg repo.CreateFeedParams) (repo.Feed, error) {
	return r.q.CreateFeed(ctx, arg)
}

func (r *feedRepo) GetFeedByID(ctx context.Context, id int32) (repo.Feed, error) {
	return r.q.GetFeedByID(ctx, id)
}

func (r *feedRepo) GetAllFeeds(ctx context.Context) ([]repo.Feed, error) {
	return r.q.GetAllFeeds(ctx)
}

func (r *feedRepo) GetFeedsByUserID(ctx context.Context, userID int32) ([]repo.Feed, error) {
	return r.q.GetFeedsByUserID(ctx, userID)
}

func (r *feedRepo) GetFeedURLByID(ctx context.Context, id int32) (string, error) {
	return r.q.GetFeedURLByID(ctx, id)
}

func (r *feedRepo) UpdateFeed(ctx context.Context, arg repo.UpdateFeedParams) (repo.Feed, error) {
	return r.q.UpdateFeed(ctx, arg)
}

func (r *feedRepo) UpdateFeedFetchTime(ctx context.Context, id int32, lastFetchedAt pgtype.Timestamptz) (repo.Feed, error) {
	return r.q.UpdateFeedFetchTime(ctx, repo.UpdateFeedFetchTimeParams{
		ID:            id,
		LastFetchedAt: lastFetchedAt,
	})
}

func (r *feedRepo) DeleteFeedByID(ctx context.Context, id int32) error {
	return r.q.DeleteFeedByID(ctx, id)
}

func (r *feedRepo) DeleteFeedIfUnused(ctx context.Context, id int32) error {
	return r.q.DeleteFeedIfUnused(ctx, id)
}

func (r *feedRepo) DeleteAllFeedsByUserID(ctx context.Context, userID int32) error {
	return r.q.DeleteAllFeedsByUserID(ctx, userID)
}
