package caching_pokemon

import (
	"github.com/tlwr/truelayer-take-home-pokemon-api/internal/pokemon"
)

type cachingPokemonClient struct {
	wrapped         pokemon.PokemonClient
	getPokemonCache map[string]*pokemon.GetPokemonResponse
}

func NewClient(c pokemon.PokemonClient) pokemon.PokemonClient {
	return &cachingPokemonClient{
		wrapped:         c,
		getPokemonCache: make(map[string]*pokemon.GetPokemonResponse, 151),
	}
}

func (c *cachingPokemonClient) Get(name string) (*pokemon.GetPokemonResponse, error) {
	if cached, found := c.getPokemonCache[name]; found && cached != nil {
		return cached, nil
	}

	cacheable, err := c.wrapped.Get(name)
	if err != nil {
		return nil, err
	}

	c.getPokemonCache[name] = cacheable
	return cacheable, nil
}
