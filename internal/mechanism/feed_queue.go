package mechanism

import "github.com/khaingminhtun/rssagg/internal/adapters/database/repo"

type FeedQueue struct {
	ch chan repo.Feed
}

func NewFeedQueue(size int) *FeedQueue {
	return &FeedQueue{
		ch: make(chan repo.Feed, size),
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

func (q *FeedQueue) Dequeue() repo.Feed {
	return <-q.ch
}

func (q *FeedQueue) Channel() <-chan repo.Feed {
	return q.ch
}
