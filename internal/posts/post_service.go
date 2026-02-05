package posts

import (
	"context"

	"github.com/khaingminhtun/rssagg/internal/adapters/database/repo"
	"github.com/khaingminhtun/rssagg/internal/pkg/errorHandle"
)

type PostService interface {
	// Get methods
	GetAllPosts(ctx context.Context) ([]PostResponse, error)
	GetPostsByFeedID(ctx context.Context, feedID int32) ([]PostResponse, error)
	GetPostByID(ctx context.Context, id int32) (*PostResponse, error)

	GetLatestPosts(ctx context.Context, limit int32) ([]PostResponse, error)
	SearchPosts(
		ctx context.Context,
		query string,
		page int32,
		limit int32,
	) ([]PostResponse, error)

	GetPostsForUser(ctx context.Context, userID int32, page int32, limit int32) ([]PostResponse, error)
}

type service struct {
	postRepo PostRepository
}

func NewPostService(postRepo PostRepository) PostService {
	return &service{
		postRepo: postRepo,
	}
}

// GetAllPosts retrieves all posts.
func (s *service) GetAllPosts(ctx context.Context) ([]PostResponse, error) {
	posts, err := s.postRepo.GetAllPosts(ctx)
	if err != nil {
		return nil, errorHandle.NotFound("no posts found")
	}

	res := make([]PostResponse, 0, len(posts))
	for _, post := range posts {
		res = append(res, mapPostToResponse(post))
	}

	return res, nil
}

// get post by ID
func (s *service) GetPostByID(ctx context.Context, id int32) (*PostResponse, error) {
	post, err := s.postRepo.GetPostByID(ctx, id)
	if err != nil {
		return nil, errorHandle.NotFound("post is not found with this id ")
	}
	res := mapPostToResponse(post)
	return &res, nil
}

// get  posts by feedID
func (s *service) GetPostsByFeedID(ctx context.Context, feedID int32) ([]PostResponse, error) {
	posts, err := s.postRepo.GetPostsByFeedID(ctx, feedID)
	if err != nil {
		return nil, errorHandle.NotFound("posts are not found with this feedid")
	}
	res := make([]PostResponse, 0, len(posts))
	for _, post := range posts {
		res = append(res, mapPostToResponse(post))
	}

	return res, nil

}

// get latest posts with limit
func (s *service) GetLatestPosts(
	ctx context.Context,
	limit int32,
) ([]PostResponse, error) {

	if limit <= 0 {
		limit = 10
	}

	posts, err := s.postRepo.GetLatestPosts(ctx, limit)
	if err != nil {
		return nil, errorHandle.NotFound("no latest posts found")
	}

	res := make([]PostResponse, 0, len(posts))
	for _, post := range posts {
		res = append(res, mapPostToResponse(post))
	}

	return res, nil
}

// search posts by query in title or description with pagination
func (s *service) SearchPosts(
	ctx context.Context,
	query string,
	page int32,
	limit int32,
) ([]PostResponse, error) {

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	offset := (page - 1) * limit

	posts, err := s.postRepo.SearchPosts(ctx, query, limit, offset)
	if err != nil {
		return nil, errorHandle.NotFound("no posts match this query")
	}

	res := make([]PostResponse, 0, len(posts))
	for _, post := range posts {
		res = append(res, mapPostToResponse(post))
	}

	return res, nil
}

// get posts for user with pagination
func (s *service) GetPostsForUser(ctx context.Context, userID int32, page int32, limit int32) ([]PostResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	offset := (page - 1) * limit

	posts, err := s.postRepo.GetPostsForUser(ctx, userID, limit, offset)
	if err != nil {
		return nil, errorHandle.NotFound("no posts found for this user")
	}

	res := make([]PostResponse, 0, len(posts))
	for _, post := range posts {
		res = append(res, mapPostToResponse(post))
	}

	return res, nil
}

func mapPostToResponse(post repo.Post) PostResponse {
	const layout = "2006-01-02 15:04:05" // custom format

	return PostResponse{
		ID:          post.ID,
		FeedID:      post.FeedID,
		Title:       post.Title,
		URL:         post.Url,
		Description: post.Description.String,
		PublishedAt: post.PublishedAt.Format(layout),
		Guid:        post.Guid,
		CreatedAt:   post.CreatedAt.Format(layout),
	}
}
