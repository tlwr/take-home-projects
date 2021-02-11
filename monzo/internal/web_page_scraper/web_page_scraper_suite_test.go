package web_page_scraper_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestWebPageScraper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "WebPageScraper Suite")
}
