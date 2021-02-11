package url_dedup_queue_test

import (
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	queue "github.com/tlwr/monzo-take-home-crawler/internal/url_dedup_queue"
)

var _ = Describe("Queue", func() {

	DescribeTable("Iter",
		func(enqueued []*url.URL, expected []*url.URL) {
			q := queue.New()

			for _, u := range enqueued {
				q.Enqueue(u)
			}

			worked := []*url.URL{}
			go func() {
				q.Iter(func(u *url.URL) {
					worked = append(worked, u)
				})
			}()

			q.Wait()

			Expect(worked).To(Equal(expected))
		},

		Entry(
			"a single URL",
			[]*url.URL{{Path: "/a"}},
			[]*url.URL{{Path: "/a"}},
		),

		Entry(
			"two distinct URLs",
			[]*url.URL{{Path: "/a"}, {Path: "/b"}},
			[]*url.URL{{Path: "/a"}, {Path: "/b"}},
		),

		Entry(
			"two identical URLs",
			[]*url.URL{{Path: "/a"}, {Path: "/a"}},
			[]*url.URL{{Path: "/a"}},
		),
	)
})
