package url_dedup_queue

import (
	"net/url"
	"sync"

	bloomfilter "github.com/kkdai/bloomfilter" // probablistic visited link checking
)

// URLDedupQueue is a queue of URLs
//
// URLs can be worked upon by multiple coroutines using the Iter function
//
// If the same URL has been enqueued multiple times, it is only worked upon once
type URLDedupQueue interface {
	// Enqueue adds a URL to the queue, if it has not been worked upon already
	Enqueue(*url.URL)

	// Iter calls the function param for each item in the queue
	// Iter is backed by a channel so can be used from multiple coroutines
	Iter(func(*url.URL))

	// Wait blocks until the queue is empty
	Wait()
}

type urlDedupQueue struct {
	c chan *url.URL

	// wg represents the capacity of the queue
	// wg.Done means the queue is empty
	wg sync.WaitGroup

	// seen will be used by multiple coroutines so must be guarded by a mutex
	// using seen is a blocking operation.
	// this represents a potential performance problem
	// but network IO elsewhere is more likely to be a bottleneck
	seen   *bloomfilter.CBF
	seenMu sync.Mutex
}

func New() URLDedupQueue {
	return &urlDedupQueue{
		c: make(chan *url.URL, 1024*1024), // should be pretty large to avoid queue saturation for pages with many links

		// expect to see 1k links with 99.95% false positive
		seen:   bloomfilter.NewCountingBloomFilter(1024, 0.05),
		seenMu: sync.Mutex{},
	}
}

func (q *urlDedupQueue) Enqueue(u *url.URL) {
	q.wg.Add(1)
	q.c <- u
}

func (q *urlDedupQueue) Iter(work func(*url.URL)) {

	for u := range q.c {
		q.seenMu.Lock()
		// Lock done in a loop so cannot use defer

		link := []byte((*u).String())

		if q.seen.Test(link) {
			q.seenMu.Unlock()
			q.wg.Done()
			continue // skip if host already seen
		} else {
			q.seen.Add(link)
			q.seenMu.Unlock()
		}

		work(u)
		q.wg.Done()
	}
}

func (q *urlDedupQueue) Wait() {
	q.wg.Wait()
	close(q.c)
}
