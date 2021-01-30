package ye_olde_pokemon

import (
	"github.com/tlwr/truelayer-take-home-pokemon-api/internal/pokemon"
	"github.com/tlwr/truelayer-take-home-pokemon-api/internal/shakespeare"
)

type yeOldePokemonClient struct {
	wrapped    pokemon.PokemonClient
	translator shakespeare.ShakespeareClient
}

func NewClient(
	c pokemon.PokemonClient,
	t shakespeare.ShakespeareClient,
) pokemon.PokemonClient {
	return &yeOldePokemonClient{
		wrapped:    c,
		translator: t,
	}
}

func (c *yeOldePokemonClient) Get(name string) (pokemon.GetPokemonResponse, error) {
	resp, err := c.wrapped.Get(name)
	if err != nil {
		return pokemon.GetPokemonResponse{}, err
	}

	translation, err := c.translator.Translate(resp.Description)
	if err != nil {
		return pokemon.GetPokemonResponse{}, err
	}

	resp.Description = translation
	return resp, nil
}
