package pokemon_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/ghttp"

	"github.com/tlwr/truelayer-take-home-pokemon-api/internal/pokemon"
)

var _ = Describe("PokemonClient", func() {
	var (
		c pokemon.PokemonClient
		s *ghttp.Server
	)

	BeforeEach(func() {
		s = ghttp.NewServer()

		c = pokemon.NewClient(s.URL(), http.DefaultClient)
	})

	AfterEach(func() {
		s.Close()
	})

	Context("happy path", func() {
		BeforeEach(func() {
			s.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/v2/pokemon-species/charizard"),
				ghttp.VerifyContentType("application/json"),
				ghttp.RespondWithJSONEncoded(200, map[string]interface{}{
					"flavor_text_entries": []map[string]interface{}{
						{
							"flavor_text": "a scary winged fire lizard",
							"language": map[string]string{
								"name": "en",
							},
						},
					},
				}),
			))
		})

		It("returns a name and description", func() {
			resp, err := c.Get("charizard")
			Expect(err).NotTo(HaveOccurred())

			Expect(resp.Name).To(Equal("charizard"))
			Expect(resp.Description).To(Equal("a scary winged fire lizard"))
		})
	})

	Context("when the status code is not HTTP OK", func() {
		BeforeEach(func() {
			s.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/v2/pokemon-species/charizard"),
				ghttp.VerifyContentType("application/json"),
				ghttp.RespondWith(403, "no thank you"),
			))
		})

		It("returns an error", func() {
			_, err := c.Get("charizard")
			Expect(err).To(MatchError("expected 200 received 403"))
		})
	})

	Context("when the response is malformed json", func() {
		BeforeEach(func() {
			s.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/v2/pokemon-species/charizard"),
				ghttp.VerifyContentType("application/json"),
				ghttp.RespondWith(200, "{ this is broken json }"),
			))
		})

		It("returns an error", func() {
			_, err := c.Get("charizard")
			Expect(err).To(MatchError(ContainSubstring("invalid character")))
		})
	})

	Context("when no english description is available", func() {
		BeforeEach(func() {
			s.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/v2/pokemon-species/charizard"),
				ghttp.VerifyContentType("application/json"),
				ghttp.RespondWithJSONEncoded(200, map[string]interface{}{
					"flavor_text_entries": []map[string]interface{}{
						{
							"flavor_text": "en bevinget brann øgle jeg er redd for",
							"language": map[string]string{
								"name": "no",
							},
						},
					},
				}),
			))
		})

		It("returns an error", func() {
			_, err := c.Get("charizard")
			Expect(err).To(MatchError(ContainSubstring("could not find english language description")))
		})
	})
})
