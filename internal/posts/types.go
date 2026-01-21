package posts

import "time"

type PostResponse struct {
	ID          int32     `json:"id"`
	FeedID      int32     `json:"feedID"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"publishedAt"`
	Guid        string    `json:"guid"`
	CreatedAt   time.Time `json:"createdAt"`
}
