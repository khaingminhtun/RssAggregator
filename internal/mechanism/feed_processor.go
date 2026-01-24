package mechanism

import (
	"context"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/khaingminhtun/rssagg/internal/adapters/database/repo"
	"github.com/khaingminhtun/rssagg/internal/feeds"
	"github.com/khaingminhtun/rssagg/internal/posts"
	"golang.org/x/net/html"
)

// ----------------------------
// FeedProcessor interface
// ----------------------------
type FeedProcessor interface {
	ProcessFeed(ctx context.Context, feed repo.Feed) error
}

// ----------------------------
// Service implementation
// ----------------------------
type FeedProcessorService struct {
	fetcher  *feeds.FetcherService
	postRepo posts.PostRepository
	feedRepo feeds.FeedRepository
}

func NewFeedProcessorService(fetcher *feeds.FetcherService, postRepo posts.PostRepository, feedRepo feeds.FeedRepository) *FeedProcessorService {
	return &FeedProcessorService{
		fetcher:  fetcher,
		postRepo: postRepo,
		feedRepo: feedRepo,
	}
}

// helper to fetch full post text from URL
func fetchPostText(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("[Processor] Failed to fetch post URL %s: %v", url, err)
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[Processor] Failed to read post body %s: %v", url, err)
		return ""
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		log.Printf("[Processor] Failed to parse HTML for %s: %v", url, err)
		return ""
	}

	var text string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			text += n.Data + " "
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return strings.Join(strings.Fields(text), " ") // trim extra whitespace
}

// ----------------------------
// Process a feed
// ----------------------------
func (s *FeedProcessorService) ProcessFeed(ctx context.Context, feed repo.Feed) error {
	if feed.FeedUrl == "" {
		log.Printf("[Processor] Skipping feed %d: empty FeedUrl", feed.ID)
		return nil
	}

	parsedFeed, err := s.fetcher.FetchAndParse(feed.FeedUrl)
	if err != nil {
		return err
	}

	for _, item := range parsedFeed.Items {
		if item.Link == "" {
			continue
		}

		publishedAt := time.Now()
		if item.PublishedParsed != nil {
			publishedAt = *item.PublishedParsed
		}

		guid := item.GUID
		if guid == "" {
			guid = item.Link
		}

		desc := item.Description
		if desc == "" {
			desc = fetchPostText(item.Link) // fetch full post text if description is empty
		}

		_, err = s.postRepo.CreatePost(ctx, repo.CreatePostParams{
			FeedID:      feed.ID,
			Title:       item.Title,
			Url:         item.Link,
			Description: desc,
			PublishedAt: publishedAt,
			Guid:        guid,
		})
		if err != nil {
			log.Printf("Failed to create post %s: %v", item.Title, err)
			return err
		}
	}

	return nil
}
