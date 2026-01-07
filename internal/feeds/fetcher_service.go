package feeds

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/khaingminhtun/rssagg/internal/utils"
	"github.com/mmcdole/gofeed"
	"golang.org/x/net/html"
)

type FetcherService struct {
	HttpClient *http.Client
}

func (f *FetcherService) DiscoverFeed(siteURL string) (string, error) {
	resp, err := f.HttpClient.Get(siteURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("failed to fetch site: %s", resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	var feedURL string
	var walker func(*html.Node)

	walker = func(n *html.Node) {
		if feedURL != "" {
			return // stop once found
		}

		if n.Type == html.ElementNode && n.Data == "link" {
			var (
				isAlternate bool
				isFeed      bool
				href        string
			)

			for _, attr := range n.Attr {
				switch strings.ToLower(attr.Key) {
				case "rel":
					if strings.Contains(strings.ToLower(attr.Val), "alternate") {
						isAlternate = true
					}
				case "type":
					val := strings.ToLower(attr.Val)
					if strings.Contains(val, "rss") || strings.Contains(val, "atom") {
						isFeed = true
					}
				case "href":
					href = attr.Val
				}
			}

			if isAlternate && isFeed && href != "" {
				feedURL = utils.ResolveURL(siteURL, href)
				return
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walker(c)
		}
	}

	walker(doc)

	// HTML discovery success
	if feedURL != "" {
		return feedURL, nil
	}

	// -------- Fallback heuristics --------

	fallbacks := []string{
		"/feed",
		"/rss",
		"/rss.xml",
		"/feed.xml",
		"/atom.xml",
		"/blog/feed",
		"/blog/feed.atom",
	}

	base, err := url.Parse(siteURL)
	if err != nil {
		return "", err
	}

	for _, path := range fallbacks {
		u := base.ResolveReference(&url.URL{Path: path})
		resp, err := f.HttpClient.Get(u.String())
		if err != nil {
			continue
		}
		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return u.String(), nil
		}
	}

	return "", errors.New("no RSS or Atom feed found")
}

func (f *FetcherService) FetchAndParse(feedURL string) (*gofeed.Feed, error) {
	parser := gofeed.NewParser()
	return parser.ParseURL(feedURL)
}
