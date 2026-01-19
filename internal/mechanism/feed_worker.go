package mechanism

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/khaingminhtun/rssagg/internal/adapters/database/repo"
)

type FeedWorker struct {
	svc *FeedProcessorService
	q   *FeedQueue
}

func NewFeedWorker(svc *FeedProcessorService, q *FeedQueue) *FeedWorker {
	return &FeedWorker{
		svc: svc,
		q:   q,
	}
}

// Start runs `n` workers concurrently
func (w *FeedWorker) Start(ctx context.Context, n int) {
	for i := 0; i < n; i++ {
		go w.run(ctx, i)
	}
}

// run is a single worker goroutine
func (w *FeedWorker) run(ctx context.Context, workerID int) {
	for {
		select {
		case feed := <-w.q.Channel():
			w.handle(feed, workerID)
		case <-ctx.Done():
			return
		}
	}
}

// handle processes a single feed and updates the fetch result
func (w *FeedWorker) handle(feed repo.Feed, workerID int) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := w.svc.ProcessFeed(ctx, feed)
	if err != nil {
		// optional logging
		log.Printf("[Worker %d] Failed to process feed %d: %v", workerID, feed.ID, err)
	}

	// Update the feed fetch result in DB
	_, updateErr := w.svc.feedRepo.UpdateFeed(ctx, repo.UpdateFeedParams{
		ID: feed.ID,
		LastFetchedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	})

	if updateErr != nil {
		// optional logging
		log.Printf("[Worker %d] Failed to update fetch result for feed %d: %v", workerID, feed.ID, updateErr)
	}
}
