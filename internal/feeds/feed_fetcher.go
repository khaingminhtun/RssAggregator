package feeds

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/khaingminhtun/rssagg/internal/pkg/utils"
	"github.com/mmcdole/gofeed"
	"golang.org/x/net/html"
)

type FetcherService struct {
	httpClient *http.Client
	sem        chan struct{}
	retry      int
}

// NewFetcherService constructs a new fetcher with timeout, concurrency limit, and retries
func NewFetcherService(timeout time.Duration, maxConcurrent int, retries int) *FetcherService {
	return &FetcherService{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		sem:   make(chan struct{}, maxConcurrent),
		retry: retries,
	}
}

// ---------------- Internal Helpers ----------------

// rate-limited, context-aware HTTP execution
func (f *FetcherService) doRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	select {
	case f.sem <- struct{}{}:
		defer func() { <-f.sem }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	return f.httpClient.Do(req)
}

// retry with linear backoff
func (f *FetcherService) doWithRetry(ctx context.Context, fn func() error) error {
	var err error
	for i := 0; i <= f.retry; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		select {
		case <-time.After(time.Duration(i+1) * time.Second):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return err
}

// ---------------- Public API ----------------

// DiscoverFeed finds the canonical RSS/Atom feed URL from a website
func (f *FetcherService) DiscoverFeed(ctx context.Context, siteURL string) (string, error) {
	var resp *http.Response

	// Fetch the page
	err := f.doWithRetry(ctx, func() error {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, siteURL, nil)
		if err != nil {
			return err
		}

		// Important headers to bypass simple WAF checks
		req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; RSSAgg/1.0; +https://yourdomain.com)")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		req.Header.Set("Connection", "keep-alive")

		resp, err = f.doRequest(ctx, req)
		return err
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return "", fmt.Errorf("failed to fetch site: %s", resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	// Look for <link rel="alternate" type="application/rss+xml">
	var feedURL string
	var walker func(*html.Node)
	walker = func(n *html.Node) {
		if feedURL != "" {
			return
		}
		if n.Type == html.ElementNode && n.Data == "link" {
			var rel, typ, href string
			for _, attr := range n.Attr {
				switch strings.ToLower(attr.Key) {
				case "rel":
					rel = strings.ToLower(attr.Val)
				case "type":
					typ = strings.ToLower(attr.Val)
				case "href":
					href = attr.Val
				}
			}
			if strings.Contains(rel, "alternate") &&
				(strings.Contains(typ, "rss") || strings.Contains(typ, "atom")) &&
				href != "" {
				feedURL = utils.ResolveURL(siteURL, href)
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walker(c)
		}
	}
	walker(doc)

	if feedURL != "" {
		return feedURL, nil
	}

	// Fallback heuristics (avoid /feed for WAF-protected sites)
	fallbacks := []string{
		"/rss.xml",
		"/feed.xml",
		"/atom.xml",
		"/blog/rss.xml",
		"/blog/feed.atom",
	}
	base, err := url.Parse(siteURL)
	if err != nil {
		return "", err
	}
	for _, path := range fallbacks {
		u := base.ResolveReference(&url.URL{Path: path})
		err := f.doWithRetry(ctx, func() error {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
			if err != nil {
				return err
			}
			req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; RSSAgg/1.0; +https://yourdomain.com)")
			resp, err = f.doRequest(ctx, req)
			return err
		})
		if err != nil {
			continue
		}
		resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 400 {
			return u.String(), nil
		}
	}

	return "", errors.New("no canonical RSS or Atom feed found")
}

// FetchAndParse fetches the feed reliably and parses RSS/Atom
func (f *FetcherService) FetchAndParse(ctx context.Context, feedURL string) (*gofeed.Feed, error) {
	var feed *gofeed.Feed
	var resp *http.Response

	err := f.doWithRetry(ctx, func() error {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
		if err != nil {
			return err
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; RSSAgg/1.0; +https://yourdomain.com)")
		req.Header.Set("Accept", "application/rss+xml, application/atom+xml, text/xml, */*")
		req.Header.Set("Connection", "keep-alive")

		resp, err = f.doRequest(ctx, req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 400 {
			return fmt.Errorf("failed to fetch feed: %s", resp.Status)
		}

		// Read body into memory
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		body = bytes.TrimSpace(body)
		if len(body) == 0 {
			return errors.New("empty feed body")
		}

		// Parse feed
		parser := gofeed.NewParser()
		feed, err = parser.ParseString(string(body))
		return err
	})

	return feed, err
}
