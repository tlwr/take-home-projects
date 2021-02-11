package ye_olde_pokemon_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	fake "github.com/tlwr/take-home-projects/truelayer/internal/fake_pokemon"
	shake "github.com/tlwr/take-home-projects/truelayer/internal/fake_shakespeare"
	"github.com/tlwr/take-home-projects/truelayer/internal/pokemon"
	yeo "github.com/tlwr/take-home-projects/truelayer/internal/ye_olde_pokemon"
)

var _ = Describe("YeOldePokemon", func() {
	var (
		f fake.FakePokemonClient
		s shake.FakeShakespeareClient
		c pokemon.PokemonClient
	)

	BeforeEach(func() {
		f = fake.FakePokemonClient{}
		s = shake.FakeShakespeareClient{}

		c = yeo.NewClient(&f, &s)
	})

	Context("happy path", func() {
		BeforeEach(func() {
			f.GetReturns(&pokemon.GetPokemonResponse{
				Name:        "pikachu",
				Description: "get the multimeter",
			}, nil)

			s.TranslateReturns("verily, a sparky boi", nil)
		})

		It("returns a translated pokemon", func() {
			resp, err := c.Get("pikachu")
			Expect(err).NotTo(HaveOccurred())

			Expect(resp.Name).To(Equal("pikachu"))
			Expect(resp.Description).To(Equal("verily, a sparky boi"))

			Expect(f.GetCallCount()).To(Equal(1))
			Expect(s.TranslateCallCount()).To(Equal(1))
		})
	})

	Context("when an the pokemon API returns an error", func() {
		BeforeEach(func() {
			f.GetReturns(&pokemon.GetPokemonResponse{}, fmt.Errorf("this is an arbitrary error"))
		})

		It("returns an error and does not attempt to translate", func() {
			_, err := c.Get("pikachu")
			Expect(err).To(MatchError("this is an arbitrary error"))

			Expect(s.TranslateCallCount()).To(Equal(0))
		})
	})

	Context("when the translator returns an error", func() {
		BeforeEach(func() {
			f.GetReturns(&pokemon.GetPokemonResponse{
				Name:        "pikachu",
				Description: "get the multimeter",
			}, nil)

			s.TranslateReturns("", fmt.Errorf("danger high voltage"))
		})

		It("returns an error and does not attempt to translate", func() {
			_, err := c.Get("pikachu")
			Expect(err).To(MatchError("danger high voltage"))
		})
	})
})
