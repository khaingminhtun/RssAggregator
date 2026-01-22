package posts

import (
	"context"

	"github.com/khaingminhtun/rssagg/internal/adapters/database/repo"
)

type PostService interface {
	// Get methods
	GetAllPosts(ctx context.Context) ([]PostResponse, error)
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
		return nil, err
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
		Guid:        post.Guid.String,
		CreatedAt:   post.CreatedAt.Format(layout),
	}
}
