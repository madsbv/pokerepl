package pokeapi

import (
	"fmt"
	"github.com/madsbv/pokerepl/internal/pokecache"
)

func GetPokemonDetails(query string, cache pokecache.Cache) (PokeapiPokemon, error) {
	url := pokemonURL(query)
	return getParsedResponse[PokeapiPokemon](url, cache)
}

var pokeapiPokemonURL string = "https://pokeapi.co/api/v2/pokemon/"

func pokemonURL(query string) string {
	return fmt.Sprintf("%s%s", pokeapiPokemonURL, query)
}
