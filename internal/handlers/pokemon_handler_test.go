package handlers_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	pokefake "github.com/tlwr/truelayer-take-home-pokemon-api/internal/fake_pokemon"
	"github.com/tlwr/truelayer-take-home-pokemon-api/internal/handlers"
	"github.com/tlwr/truelayer-take-home-pokemon-api/internal/pokemon"
)

var _ = Describe("PokemonHandler", func() {
	var (
		req      *http.Request
		recorder *httptest.ResponseRecorder

		h  *handlers.PokemonHandler
		pc *pokefake.FakePokemonClient
	)

	BeforeEach(func() {
		req = &http.Request{}
		recorder = httptest.NewRecorder()

		pc = &pokefake.FakePokemonClient{}
		h = handlers.NewPokemonHandler(pc)
	})

	Describe("Get", func() {
		Context("happy path", func() {
			BeforeEach(func() {
				pc.GetReturns(&pokemon.GetPokemonResponse{
					Name:        "pikachu",
					Description: "a sparky boi",
				}, nil)

				req.URL = &(url.URL{Path: "/pokemon/pikachu"})
			})

			It("returns 200 with name and description", func() {
				h.HandleGet(recorder, req)

				resp := recorder.Result()

				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(200))
				Expect(resp.Header.Get("Content-Type")).To(Equal("application/json"))

				b, err := ioutil.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(b).To(MatchJSON(`{"name": "pikachu", "description": "a sparky boi"}`))
			})
		})

		Context("when the pokemon client returns an error", func() {
			BeforeEach(func() {
				pc.GetReturns(&pokemon.GetPokemonResponse{}, fmt.Errorf("oh no"))

				req.URL = &(url.URL{Path: "/pokemon/pikachu"})
			})

			It("returns 500 with name and description", func() {
				h.HandleGet(recorder, req)

				resp := recorder.Result()

				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(500))
				Expect(resp.Header.Get("Content-Type")).To(Equal("application/json"))

				b, err := ioutil.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(b).To(MatchJSON(`{"message": "oh no"}`))
			})
		})

		Context("when the pokemon client does not return a pokemon", func() {
			BeforeEach(func() {
				pc.GetReturns(nil, nil)

				req.URL = &(url.URL{Path: "/pokemon/pikachu"})
			})

			It("returns 404 with a descriptive message", func() {
				h.HandleGet(recorder, req)

				resp := recorder.Result()

				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(404))
				Expect(resp.Header.Get("Content-Type")).To(Equal("application/json"))

				b, err := ioutil.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(b).To(MatchJSON(`{"message": "not found"}`))
			})
		})
	})
})
