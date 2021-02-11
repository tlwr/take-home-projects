package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/tlwr/take-home-projects/truelayer/internal/pokemon"
)

var (
	pokemonNameValidator = regexp.MustCompile("^[a-z]+$")
)

type PokemonHandler struct {
	pc pokemon.PokemonClient
}

type GetPokemonResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GetPokemonErrorResponse struct {
	Message string `json:"message"`
}

func writeJSON(rw http.ResponseWriter, status int, o interface{}) {
	if b, err := json.Marshal(o); err != nil {
		rw.WriteHeader(500)
		rw.Write([]byte("could not marshal json"))
	} else {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(status)
		rw.Write(b)
	}
}

func (h *PokemonHandler) HandleGet(rw http.ResponseWriter, req *http.Request) {
	pokemonName := strings.TrimPrefix(req.URL.Path, "/pokemon/")

	if !pokemonNameValidator.MatchString(pokemonName) {
		writeJSON(rw, 400, GetPokemonErrorResponse{Message: "could not validate pokemon name"})
		return
	}

	resp, err := h.pc.Get(strings.ToLower(pokemonName))
	if err != nil {
		writeJSON(rw, 500, GetPokemonErrorResponse{Message: fmt.Sprintf("%s", err)})
		return
	}

	if resp == nil {
		writeJSON(rw, 404, GetPokemonErrorResponse{Message: "not found"})
		return
	}

	writeJSON(rw, 200, GetPokemonResponse{
		Name:        resp.Name,
		Description: resp.Description,
	})
}

func NewPokemonHandler(pc pokemon.PokemonClient) *PokemonHandler {
	return &PokemonHandler{
		pc: pc,
	}
}
