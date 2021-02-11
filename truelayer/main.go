package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	caching "github.com/tlwr/truelayer-take-home-pokemon-api/internal/caching_pokemon"
	"github.com/tlwr/truelayer-take-home-pokemon-api/internal/pokemon"
	"github.com/tlwr/truelayer-take-home-pokemon-api/internal/shakespeare"
	yeo "github.com/tlwr/truelayer-take-home-pokemon-api/internal/ye_olde_pokemon"

	"github.com/tlwr/truelayer-take-home-pokemon-api/internal/handlers"

	pokefake "github.com/tlwr/truelayer-take-home-pokemon-api/internal/fake_pokemon"
	fakespeare "github.com/tlwr/truelayer-take-home-pokemon-api/internal/fake_shakespeare"
)

func main() {
	var (
		pokemonAPI     = "https://pokeapi.co"
		shakespeareAPI = "https://api.funtranslations.com"

		pokemonClient     pokemon.PokemonClient
		shakespeareClient shakespeare.ShakespeareClient

		stubAPIs bool

		bind string
		port int
	)

	// we want to stub APIs because the shakespeare API allegedly has draconian rate limit
	flag.StringVar(&bind, "bind", "", "interface or address to bind to, default all")
	flag.BoolVar(&stubAPIs, "stubs", false, "if true, use stub APIs")
	flag.IntVar(&port, "port", 5000, "port on which server should run")
	flag.Parse()

	if port <= 1024 || port > 65535 {
		log.Fatalf("port (%d) should be unprivileged, ie between >1024 and <=65535", port)
	}

	// raw API clients
	if stubAPIs {
		log.Println("stubbing apis")

		fakePokemonClient := &pokefake.FakePokemonClient{}
		fakePokemonClient.GetCalls(func(name string) (*pokemon.GetPokemonResponse, error) {
			if name == "pikachu" || name == "charizard" {
				return &pokemon.GetPokemonResponse{Name: name, Description: ""}, nil
			} else if name == "missingno" {
				return nil, fmt.Errorf("missingno")
			} else {
				return nil, nil
			}
		})
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

	handler := handlers.NewPokemonHandler(pokemonClient)

	// server
	mux := http.NewServeMux()
	mux.HandleFunc("/pokemon/", handler.HandleGet)

	s := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", bind, port),
		Handler: mux,
	}

	log.Printf("listening on %s:%d", bind, port)
	log.Fatal(s.ListenAndServe())
}
