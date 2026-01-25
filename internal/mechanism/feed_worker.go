package mechanism

import (
	"context"
	"log"
	"time"

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
		feed, ok := w.q.Dequeue(ctx)
		if !ok {
			return
		}

		w.handle(ctx, feed, workerID)
	}
}

// handle processes a single feed and updates the fetch result
func (w *FeedWorker) handle(
	parentCtx context.Context,
	feed repo.Feed,
	workerID int,
) {
	ctx, cancel := context.WithTimeout(parentCtx, 30*time.Second)
	defer cancel()

	if err := w.svc.ProcessFeed(ctx, feed); err != nil {
		log.Printf("[Worker %d] feed=%d error=%v", workerID, feed.ID, err)
		return
	}

	if err := w.svc.MarkFetched(ctx, feed.ID); err != nil {
		log.Printf("[Worker %d] update failed feed=%d err=%v", workerID, feed.ID, err)
	}
}
