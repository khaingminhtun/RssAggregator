package feeds

type RequestSiteURL struct {
	SiteURL string `json:"siteURL"`
}

type FeedResponse struct {
	Title   string `json:"title"`
	SiteURL string `json:"siteURL"`
	FeedURL string `json:"feedURL"`
}
