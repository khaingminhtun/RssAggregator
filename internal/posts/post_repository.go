package posts

import (
	"context"

	"github.com/khaingminhtun/rssagg/internal/adapters/database/repo"
)

type PostRepository interface {
	CreatePost(ctx context.Context, arg repo.CreatePostParams) (repo.Post, error)
	GetAllPosts(ctx context.Context) ([]repo.Post, error)
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
