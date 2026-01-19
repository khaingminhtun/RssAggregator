package mechanism

import (
	"context"
	"log"
	"time"

	"github.com/khaingminhtun/rssagg/internal/feeds"
)

// FeedScheduler periodically fetches feeds and pushes them to a queue
type FeedScheduler struct {
	feedRepo feeds.FeedRepository
	queue    *FeedQueue
	interval time.Duration
}

// NewFeedScheduler creates a scheduler
func NewFeedScheduler(feedRepo feeds.FeedRepository, queue *FeedQueue, interval time.Duration) *FeedScheduler {
	return &FeedScheduler{
		feedRepo: feedRepo,
		queue:    queue,
		interval: interval,
	}
}

// Start begins the periodic scheduling
func (s *FeedScheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.scheduleFeeds(ctx)
		case <-ctx.Done():
			log.Println("FeedScheduler stopped")
			return
		}
	}
}

// scheduleFeeds fetches all feeds and pushes them into the queue
func (s *FeedScheduler) scheduleFeeds(ctx context.Context) {
	feedsList, err := s.feedRepo.GetAllFeeds(ctx)
	if err != nil {
		log.Println("FeedScheduler: failed to fetch feeds:", err)
		return
	}

	for _, f := range feedsList {
		if !s.queue.TryEnqueue(f) {
			log.Printf("FeedScheduler: queue full, skipping feed %d", f.ID)
		}
	}
}
