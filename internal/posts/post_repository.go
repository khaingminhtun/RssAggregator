package posts

import (
	"context"

	"github.com/khaingminhtun/rssagg/internal/adapters/database/repo"
)

type PostRepository interface {
	CreatePost(ctx context.Context, arg repo.CreatePostParams) (repo.Post, error)

	// Optional (recommended for CRUD completeness)
	GetPostsByFeedID(ctx context.Context, feedID int32) ([]repo.Post, error)
	DeletePostsByFeedID(ctx context.Context, feedID int32) error
}
