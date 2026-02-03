package feeds

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
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

func NewFetcherService(timeout time.Duration, maxConcurrent int, retries int) *FetcherService {
	return &FetcherService{
		httpClient: &http.Client{Timeout: timeout},
		sem:        make(chan struct{}, maxConcurrent),
		retry:      retries,
	}
}

// ---------------- Internal Helpers ----------------

// fetchURL handles HTTP GET with headers, gzip, retries, and concurrency limit
func (f *FetcherService) fetchURL(ctx context.Context, url string) ([]byte, error) {
	var body []byte
	var err error

	err = f.doWithRetry(ctx, func() error {
		select {
		case f.sem <- struct{}{}:
			defer func() { <-f.sem }()
		case <-ctx.Done():
			return ctx.Err()
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; RSSAgg/1.0)")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("Connection", "keep-alive")

		resp, err := f.httpClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 400 {
			return fmt.Errorf("failed to fetch url: %s", resp.Status)
		}

		reader := io.Reader(resp.Body)
		if resp.Header.Get("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(resp.Body)
			if err != nil {
				return err
			}
			defer gz.Close()
			reader = gz
		}

		body, err = io.ReadAll(reader)
		if err != nil {
			return err
		}
		body = bytes.TrimSpace(body)
		if len(body) == 0 {
			return errors.New("empty body")
		}
		return nil
	})

	return body, err
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

// DiscoverFeed finds the canonical RSS/Atom feed URL
func (f *FetcherService) DiscoverFeed(ctx context.Context, siteURL string) (string, error) {
	// 1. Try parsing HTML for <link rel="alternate">
	body, err := f.fetchURL(ctx, siteURL)
	if err == nil {
		doc, err := html.Parse(bytes.NewReader(body))
		if err == nil {
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
		}
	}

	// 2. Fallback URLs
	fallbacks := []string{"/feed", "/rss.xml", "/feed.xml", "/atom.xml", "/blog/rss.xml", "/blog/feed.atom"}
	base, err := url.Parse(siteURL)
	if err != nil {
		return "", err
	}

	for _, path := range fallbacks {
		u := base.ResolveReference(&url.URL{Path: path}).String()
		body, err := f.fetchURL(ctx, u)
		if err == nil && len(body) > 0 {
			return u, nil
		}
	}

	return "", errors.New("no canonical RSS or Atom feed found")
}

// FetchAndParse fetches the feed and parses it
func (f *FetcherService) FetchAndParse(ctx context.Context, feedURL string) (*gofeed.Feed, error) {
	body, err := f.fetchURL(ctx, feedURL)
	if err != nil {
		return nil, err
	}

	// Use io.Reader directly for more robust parsing
	parser := gofeed.NewParser()
	feed, err := parser.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse feed: %w", err)
	}
	return feed, nil
}
