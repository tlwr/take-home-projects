package main

import (
	"flag"
	"log"
	"net/http"

	caching "github.com/tlwr/truelayer-take-home-pokemon-api/internal/caching_pokemon"
	"github.com/tlwr/truelayer-take-home-pokemon-api/internal/pokemon"
	"github.com/tlwr/truelayer-take-home-pokemon-api/internal/shakespeare"
	yeo "github.com/tlwr/truelayer-take-home-pokemon-api/internal/ye_olde_pokemon"

	pokefake "github.com/tlwr/truelayer-take-home-pokemon-api/internal/fake_pokemon"
	fakespeare "github.com/tlwr/truelayer-take-home-pokemon-api/internal/fake_shakespeare"
)

func main() {
	var (
		testPokemon = "pikachu"

		pokemonAPI     = "https://pokeapi.co"
		shakespeareAPI = "https://api.funtranslations.com"

		pokemonClient     pokemon.PokemonClient
		shakespeareClient shakespeare.ShakespeareClient

		stubAPIs bool
	)

	// we want to stub APIs because the shakespeare API allegedly has draconian rate limit
	flag.BoolVar(&stubAPIs, "stubs", false, "if true, use stub APIs")
	flag.Parse()

	// raw API clients
	if stubAPIs {
		log.Println("stubbing apis")

		fakePokemonClient := &pokefake.FakePokemonClient{}
		fakePokemonClient.GetReturns(pokemon.GetPokemonResponse{Name: "shrew", Description: ""}, nil)
		pokemonClient = fakePokemonClient

		fakeShakespeareClient := &fakespeare.FakeShakespeareClient{}
		fakeShakespeareClient.TranslateReturns("you speak an infinite deal of nothing", nil)
		shakespeareClient = fakeShakespeareClient
	} else {
		pokemonClient = pokemon.NewClient(pokemonAPI, http.DefaultClient)
		shakespeareClient = shakespeare.NewClient(shakespeareAPI, http.DefaultClient)
	}

	// our pokemon client now translates descriptions
	pokemonClient = yeo.NewClient(pokemonClient, shakespeareClient)

	// our pokemon client only requests a pokemon once
	pokemonClient = caching.NewClient(pokemonClient)

	log.Printf("requesting test pokemon (%s)", testPokemon)
	resp, err := pokemonClient.Get(testPokemon)
	if err != nil {
		log.Fatalf("error requesting test pokemon (%s): %s", testPokemon, err)
	}
	log.Printf("retrieved test pokemon (%s)", resp.Name)
	log.Printf(`test pokemon (%s) is described as "%s"`, resp.Name, resp.Description)
}
