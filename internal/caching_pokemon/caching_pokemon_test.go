package caching_pokemon_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	caching "github.com/tlwr/truelayer-take-home-pokemon-api/internal/caching_pokemon"
	fake "github.com/tlwr/truelayer-take-home-pokemon-api/internal/fake_pokemon"
	"github.com/tlwr/truelayer-take-home-pokemon-api/internal/pokemon"
)

var _ = Describe("CachingPokemon", func() {
	var (
		f fake.FakePokemonClient
		c pokemon.PokemonClient
	)

	BeforeEach(func() {
		f = fake.FakePokemonClient{}
		c = caching.NewClient(&f)
	})

	Context("happy path", func() {
		BeforeEach(func() {
			f.GetReturns(&pokemon.GetPokemonResponse{
				Name:        "pikachu",
				Description: "get the multimeter",
			}, nil)
		})

		It("returns a pokemon", func() {
			resp, err := c.Get("pikachu")
			Expect(err).NotTo(HaveOccurred())

			Expect(resp.Name).To(Equal("pikachu"))
			Expect(resp.Description).To(Equal("get the multimeter"))
		})
	})

	Context("when called multiple times", func() {
		BeforeEach(func() {
			f.GetReturns(&pokemon.GetPokemonResponse{
				Name:        "pikachu",
				Description: "get the multimeter",
			}, nil)
		})

		It("returns a pokemon having called the wrapped client only once", func() {
			for i := 1; i <= 2; i++ {
				resp, err := c.Get("pikachu")
				Expect(err).NotTo(HaveOccurred())

				Expect(resp.Name).To(Equal("pikachu"))
				Expect(resp.Description).To(Equal("get the multimeter"))

				Expect(f.GetCallCount()).To(Equal(1))
			}
		})
	})

	Context("when an error is returned", func() {
		BeforeEach(func() {
			f.GetReturns(nil, fmt.Errorf("this is an arbitrary error"))
		})

		It("returns an error", func() {
			_, err := c.Get("pikachu")
			Expect(err).To(MatchError("this is an arbitrary error"))
		})
	})

	Context("when nil is returned", func() {
		BeforeEach(func() {
			f.GetReturns(nil, nil)
		})

		It("returns an nil", func() {
			resp, err := c.Get("pikachu")
			Expect(resp).To(BeNil())
			Expect(err).To(BeNil())
		})

		It("does not cache", func() {
			resp, err := c.Get("pikachu")
			Expect(resp).To(BeNil())
			Expect(err).To(BeNil())

			resp, err = c.Get("pikachu")
			Expect(resp).To(BeNil())
			Expect(err).To(BeNil())

			Expect(f.GetCallCount()).To(Equal(2))
		})
	})
})
