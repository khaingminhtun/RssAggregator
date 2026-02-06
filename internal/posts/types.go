package posts

type PostResponse struct {
	ID          int32  `json:"id"`
	FeedID      int32  `json:"feedID"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
	PublishedAt string `json:"publishedAt"`
	Guid        string `json:"guid"`
	CreatedAt   string `json:"createdAt"`
}

type MarkReadRequest struct {
	IsRead bool `json:"isRead"`
}

type MarkFavouriteRequest struct {
	IsFavourite bool `json:"isFavourite"`
}
