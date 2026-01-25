package mechanism

import (
	"context"

	"github.com/khaingminhtun/rssagg/internal/adapters/database/repo"
)

type FeedQueue struct {
	ch chan repo.Feed
}

func NewFeedQueue(size int) *FeedQueue {
	return &FeedQueue{
		ch: make(chan repo.Feed, size),
	}
}

func (q *FeedQueue) Enqueue(ctx context.Context, feed repo.Feed) error {
	select {
	case q.ch <- feed:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (q *FeedQueue) TryEnqueue(feed repo.Feed) bool {
	select {
	case q.ch <- feed:
		return true
	default:
		return false
	}
}

func (q *FeedQueue) Dequeue(ctx context.Context) (repo.Feed, bool) {
	select {
	case feed, ok := <-q.ch:
		return feed, ok
	case <-ctx.Done():
		return repo.Feed{}, false
	}
}
