package main

import (
	"log"
	"net/http"

	caching "github.com/tlwr/truelayer-take-home-pokemon-api/internal/caching_pokemon"
	"github.com/tlwr/truelayer-take-home-pokemon-api/internal/pokemon"
	"github.com/tlwr/truelayer-take-home-pokemon-api/internal/shakespeare"
	yeo "github.com/tlwr/truelayer-take-home-pokemon-api/internal/ye_olde_pokemon"
)

func main() {
	var (
		testPokemon = "pikachu"

		pokemonAPI     = "https://pokeapi.co"
		shakespeareAPI = "https://api.funtranslations.com"

		pokemonClient     pokemon.PokemonClient
		shakespeareClient shakespeare.ShakespeareClient
	)

	// raw API clients
	shakespeareClient = shakespeare.NewClient(shakespeareAPI, http.DefaultClient)
	pokemonClient = pokemon.NewClient(pokemonAPI, http.DefaultClient)

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
