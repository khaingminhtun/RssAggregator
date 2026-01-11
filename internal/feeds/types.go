package feeds

import "time"

type RequestSiteURL struct {
	SiteURL string `json:"siteURL"`
}

type FeedResponse struct {
	ID            int32     `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	WebsiteUrl    string    `json:"websiteUrl"`
	FeedURL       string    `json:"feedURL"`
	CreatedAt     time.Time `json:"createdAt"`
	LastFetchedAt time.Time `json:"lastFetchedAt"`
}

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
