package posts

import (
	"context"

	"github.com/khaingminhtun/rssagg/internal/adapters/database/repo"
)

type PostRepository interface {
	CreatePost(ctx context.Context, arg repo.CreatePostParams) (repo.Post, error)
	GetAllPosts(ctx context.Context) ([]repo.Post, error)
	GetPostByID(ctx context.Context, id int32) (repo.Post, error)
	GetPostsByFeedID(ctx context.Context, feedID int32) ([]repo.Post, error)
}

// --- Concrete SQLC-backed implementation --- //

type postRepo struct {
	q *repo.Queries
}

// NewPostRepository returns a PostRepository implemented with SQLC queries.
func NewPostRepository(db repo.DBTX) PostRepository {
	return &postRepo{
		q: repo.New(db),
	}
}

func (r *postRepo) CreatePost(ctx context.Context, arg repo.CreatePostParams) (repo.Post, error) {
	return r.q.CreatePost(ctx, arg)
}

// get all posts
func (r *postRepo) GetAllPosts(ctx context.Context) ([]repo.Post, error) {
	return r.q.GetAllPosts(ctx)
}

// get post by ID
func (r *postRepo) GetPostByID(ctx context.Context, id int32) (repo.Post, error) {
	return r.q.GetPostByID(ctx, id)
}

// get posts by feed ID
func (r *postRepo) GetPostsByFeedID(ctx context.Context, feedID int32) ([]repo.Post, error) {
	return r.q.GetPostsByFeedID(ctx, feedID)
}
