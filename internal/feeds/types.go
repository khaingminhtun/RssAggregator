package feeds

type RequestSiteURL struct {
	SiteURL string `json:"siteURL"`
}

type FeedResponse struct {
	ID            int32  `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	WebsiteUrl    string `json:"websiteUrl"`
	FeedURL       string `json:"feedURL"`
	CreatedAt     string `json:"createdAt"`
	LastFetchedAt string `json:"lastFetchedAt"`
}
