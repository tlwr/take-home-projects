package web_page_scraper_test

import (
	"net/http"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jarcoal/httpmock"

	scraper "github.com/tlwr/monzo-take-home-crawler/internal/web_page_scraper"
)

var _ = Describe("Scraper", func() {
	const (
		monzo = "https://www.monzo.com/start"
	)

	var (
		client  *http.Client
		current *url.URL

		s scraper.Scraper

		res *scraper.ScrapeResult

		responder httpmock.Responder
	)

	BeforeEach(func() {
		client = &http.Client{}
		httpmock.ActivateNonDefault(client)

		var err error
		current, err = url.Parse(monzo)
		Expect(err).NotTo(HaveOccurred())

		s = scraper.New(client)
	})

	AfterEach(func() {
		httpmock.DeactivateAndReset()
	})

	JustBeforeEach(func() {
		httpmock.RegisterResponder(
			"GET",
			monzo,
			responder,
		)

		var err error
		res, err = s.Scrape(current)
		Expect(err).NotTo(HaveOccurred())
		Expect(res.URL).To(Equal(current))
	})

	Context("happy path", func() {
		BeforeEach(func() {
			responder = httpmock.NewStringResponder(
				200,
				`
<html>
<body>
<a href="/host/relative"></a>
<a href="page/relative"></a>
<a href="https://web.monzo.com/absolute"></a>
</body>
</html>`,
			)
		})

		It("aggregates the links", func() {
			Expect(res.Links).To(WithTransform(linksToStrings, ConsistOf(
				"https://www.monzo.com/start/page/relative",
				"https://www.monzo.com/host/relative",
				"https://web.monzo.com/absolute",
			)))
		})
	})

	Context("when mailto links and local anchors are present", func() {
		BeforeEach(func() {
			responder = httpmock.NewStringResponder(
				200,
				`
<html>
<body>
<a href="https://foo.bar"></a>
<a href="#anchor"></a>
<a href="mailto:root@eruditorum.org"></a>
</body>
</html>`,
			)
		})

		It("discards the links", func() {
			Expect(res.Links).To(WithTransform(linksToStrings, ConsistOf(
				"https://foo.bar",
			)))
		})
	})

	Context("when links are not parseable", func() {
		BeforeEach(func() {
			responder = httpmock.NewStringResponder(
				200,
				`
<html>
<body>
<a href="http://empty-port:/bar"></a>
</body>
</html>`,
			)
		})

		It("aggregates the parse errors", func() {
			Expect(res.ParseErrors).To(ConsistOf(MatchError(ContainSubstring("empty port"))))
		})
	})
})

func linksToStrings(urls []*url.URL) []string {
	links := []string{}

	for _, u := range urls {
		Expect(u).NotTo(BeNil())
		links = append(links, u.String())
	}

	return links
}
