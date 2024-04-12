package pokeapi

import (
	"fmt"
	"github.com/madsbv/pokerepl/internal/pokecache"
)

var pokeapiLocationURL string = "https://pokeapi.co/api/v2/location-area/"

func locationURL(query string) string {
	return fmt.Sprintf("%s%s", pokeapiLocationURL, query)
}

func GetLocationDetails(query string, cache pokecache.Cache) (PokeapiLocation, error) {
	url := locationURL(query)
	return getParsedResponse[PokeapiLocation](url, cache)
}

func GetLocations(q *string, cache pokecache.Cache) (LocationList, error) {
	query := ""
	if q == nil {
		query = pokeapiLocationURL
	} else {
		query = *q
	}

	return getParsedResponse[LocationList](query, cache)
}
