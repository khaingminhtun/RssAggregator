package feeds

import (
	"errors"
	"net/http"
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

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	var feedURL string
	var walker func(*html.Node)
	walker = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "link" {
			var isRSS bool
			var href string

			for _, attr := range n.Attr {
				if attr.Key == "type" && strings.Contains(attr.Val, "rss") {
					isRSS = true
				}
				if attr.Key == "href" {
					href = attr.Val
				}
			}

			if isRSS && href != "" {
				feedURL = href
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walker(c)
		}
	}
	walker(doc)

	if feedURL == "" {
		return "", errors.New("no RSS feed found")
	}

	return utils.ResolveURL(siteURL, feedURL), nil
}

func (f *FetcherService) FetchAndParse(feedURL string) (*gofeed.Feed, error) {
	parser := gofeed.NewParser()
	return parser.ParseURL(feedURL)
}
