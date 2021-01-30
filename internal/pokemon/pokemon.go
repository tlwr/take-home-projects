package pokemon

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ../fake_pokemon/fake_pokemon.go . PokemonClient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type GetPokemonResponse struct {
	Name        string
	Description string
}

type PokemonClient interface {
	Get(name string) (GetPokemonResponse, error)
}

type pokemonClient struct {
	client  *http.Client
	baseURL string
}

type pokemonSpeciesResponse struct {
	// Fields are exported so encoding/json can unmarshal

	FlavorTextEntries []struct {
		FlavorText string `json:"flavor_text"`
		Language   struct {
			Name string `json:"name"`
		} `json:"language"`
	} `json:"flavor_text_entries"`
}

func (c *pokemonClient) Get(name string) (p GetPokemonResponse, err error) {
	url := fmt.Sprintf("%s/api/v2/pokemon-species/%s", c.baseURL, name)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("expected 200 received %d", resp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var parsed pokemonSpeciesResponse
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		return
	}

	var description string
	for _, flavorText := range parsed.FlavorTextEntries {
		if flavorText.Language.Name == "en" && len(flavorText.FlavorText) > len(description) {
			description = flavorText.FlavorText
		}
	}

	if len(description) == 0 {
		err = fmt.Errorf("could not find english language description")
		return
	}

	p.Name = name
	p.Description = description
	return
}

func NewClient(baseURL string, client *http.Client) PokemonClient {
	return &pokemonClient{
		baseURL: baseURL,
		client:  client,
	}
}
