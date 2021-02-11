package shakespeare_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/ghttp"

	"github.com/tlwr/take-home-projects/truelayer/internal/shakespeare"
)

var _ = Describe("ShakespeareClient", func() {
	var (
		c shakespeare.ShakespeareClient
		s *ghttp.Server
	)

	BeforeEach(func() {
		s = ghttp.NewServer()

		c = shakespeare.NewClient(s.URL(), http.DefaultClient)
	})

	AfterEach(func() {
		s.Close()
	})

	Context("happy path", func() {
		BeforeEach(func() {
			s.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/translate/shakespeare"),
				ghttp.VerifyFormKV("text", "humans are proud and stupid"),
				ghttp.RespondWithJSONEncoded(200, map[string]interface{}{
					"success": map[string]interface{}{
						"total": 1,
					},
					"contents": map[string]interface{}{
						"translated": "the common curse of mankind, folly and ignorance",
					},
				}),
			))
		})

		It("translates the text", func() {
			translation, err := c.Translate("humans are proud and stupid")
			Expect(err).NotTo(HaveOccurred())

			Expect(translation).To(Equal("the common curse of mankind, folly and ignorance"))
		})
	})

	Context("when the status code is not HTTP OK", func() {
		BeforeEach(func() {
			s.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/translate/shakespeare"),
				ghttp.VerifyFormKV("text", "alas poor network, i knew it well"),
				ghttp.RespondWith(403, "no thank you"),
			))
		})

		It("returns an error", func() {
			_, err := c.Translate("alas poor network, i knew it well")
			Expect(err).To(MatchError("expected 200 received 403"))
		})
	})

	Context("when the response is malformed json", func() {
		BeforeEach(func() {
			s.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/translate/shakespeare"),
				ghttp.VerifyFormKV("text", "json or not json, that is the question"),
				ghttp.RespondWith(200, "{ alas this is broken json }"),
			))
		})

		It("returns an error", func() {
			_, err := c.Translate("json or not json, that is the question")
			Expect(err).To(MatchError(ContainSubstring("invalid character")))
		})
	})
})
