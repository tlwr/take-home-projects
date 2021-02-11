package ye_olde_pokemon

import (
	"github.com/tlwr/take-home-projects/truelayer/internal/pokemon"
	"github.com/tlwr/take-home-projects/truelayer/internal/shakespeare"
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

func (c *yeOldePokemonClient) Get(name string) (*pokemon.GetPokemonResponse, error) {
	resp, err := c.wrapped.Get(name)
	if resp == nil || err != nil {
		return nil, err
	}

	translation, err := c.translator.Translate(resp.Description)
	if err != nil {
		return nil, err
	}

	resp.Description = translation
	return resp, nil
}
