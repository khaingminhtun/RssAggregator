package posts

import (
	"context"

	"github.com/khaingminhtun/rssagg/internal/adapters/database/repo"
	"github.com/khaingminhtun/rssagg/internal/pkg/utils"
)

type PostRepository interface {
	CreatePost(ctx context.Context, arg repo.CreatePostParams) (repo.Post, error)
	GetAllPosts(ctx context.Context) ([]repo.Post, error)
	GetPostByID(ctx context.Context, id int32) (repo.Post, error)
	GetPostsByFeedID(ctx context.Context, feedID int32) ([]repo.Post, error)

	GetLatestPosts(ctx context.Context, limit int32) ([]repo.Post, error)
	SearchPosts(ctx context.Context, query string, limit int32, offset int32) ([]repo.Post, error)
	GetPostsForUser(ctx context.Context, userID int32, limit int32, offset int32) ([]repo.Post, error)

	MarkPostRead(ctx context.Context, userID int32, postID int32, isRead bool) error
	MarkPostFavourite(ctx context.Context, userID int32, postID int32, isFavourite bool) error
	GetFavouritePostsForUser(ctx context.Context, userID int32, limit int32, offset int32) ([]repo.Post, error)
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

// get latest posts with limit
func (r *postRepo) GetLatestPosts(ctx context.Context, limit int32) ([]repo.Post, error) {
	return r.q.GetLatestPosts(ctx, limit)
}

// search posts by query in title or description
func (r *postRepo) SearchPosts(ctx context.Context, query string, limit int32, offset int32) ([]repo.Post, error) {
	return r.q.SearchPosts(ctx, repo.SearchPostsParams{
		Column1: utils.NullString("%" + query + "%"),
		Limit:   limit,
		Offset:  offset,
	})
}

// get posts for user with pagination
func (r *postRepo) GetPostsForUser(ctx context.Context, userID int32, limit int32, offset int32) ([]repo.Post, error) {
	return r.q.GetPostsForUser(ctx, repo.GetPostsForUserParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
}

// mark post as read/unread for user
func (r *postRepo) MarkPostRead(ctx context.Context, userID int32, postID int32, isRead bool) error {
	return r.q.MarkPostRead(ctx, repo.MarkPostReadParams{
		UserID: userID,
		PostID: postID,
		IsRead: isRead,
	})
}

// mark post as favourite/unfavourite for user
func (r *postRepo) MarkPostFavourite(ctx context.Context, userID int32, postID int32, isFavourite bool) error {
	return r.q.MarkPostFavorite(ctx, repo.MarkPostFavoriteParams{
		UserID:     userID,
		PostID:     postID,
		IsFavorite: isFavourite,
	})
}

// get favourite posts for user with pagination
func (r *postRepo) GetFavouritePostsForUser(ctx context.Context, userID int32, limit int32, offset int32) ([]repo.Post, error) {
	return r.q.GetFavoritePostsByUser(ctx, repo.GetFavoritePostsByUserParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
}
