package link_parser_test

import (
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/tlwr/monzo-take-home-crawler/internal/link_parser"
)

var _ = Describe("Parser", func() {
	var (
		current *url.URL
	)

	BeforeEach(func() {
		var err error
		current, err = url.Parse("https://www.monzo.com/hot/coral")
		Expect(err).NotTo(HaveOccurred())
	})

	It("handles relative links", func() {
		link, err := Parse(current, "card")
		Expect(err).NotTo(HaveOccurred())

		Expect(link.Scheme).To(Equal("https"))
		Expect(link.Host).To(Equal("www.monzo.com"))
		Expect(link.Path).To(Equal("/hot/coral/card"))
	})

	It("handles host relative links", func() {
		link, err := Parse(current, "/premium")
		Expect(err).NotTo(HaveOccurred())

		Expect(link.Scheme).To(Equal("https"))
		Expect(link.Host).To(Equal("www.monzo.com"))
		Expect(link.Path).To(Equal("/premium"))
	})

	It("handles absolute http links", func() {
		link, err := Parse(current, "http://badssl.com")
		Expect(err).NotTo(HaveOccurred())

		Expect(link.Scheme).To(Equal("http"))
		Expect(link.Host).To(Equal("badssl.com"))
		Expect(link.Path).To(Equal(""))
	})

	It("handles absolute https links", func() {
		link, err := Parse(current, "https://web.monzo.com/begin")
		Expect(err).NotTo(HaveOccurred())

		Expect(link.Scheme).To(Equal("https"))
		Expect(link.Host).To(Equal("web.monzo.com"))
		Expect(link.Path).To(Equal("/begin"))
	})

	It("handles absolute https links with anchors", func() {
		link, err := Parse(current, "https://web.monzo.com/begin#footer")
		Expect(err).NotTo(HaveOccurred())

		Expect(link.Scheme).To(Equal("https"))
		Expect(link.Host).To(Equal("web.monzo.com"))
		Expect(link.Path).To(Equal("/begin"))
	})

	It("discards mailto: links", func() {
		link, err := Parse(current, "mailto:root@eruditorum.org")
		Expect(link).To(BeNil())
		Expect(err).To(BeNil())
	})

	It("discards current page anchors", func() {
		link, err := Parse(current, "#footer")
		Expect(link).To(BeNil())
		Expect(err).To(BeNil())
	})
})
