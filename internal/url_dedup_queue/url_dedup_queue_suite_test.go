package url_dedup_queue_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestUrlDedupQueue(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UrlDedupQueue Suite")
}
