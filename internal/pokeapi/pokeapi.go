package pokeapi

import (
	"encoding/json"
	"fmt"
	"github.com/madsbv/pokerepl/internal/pokecache"
	"io"
	"net/http"
)

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

func getParsedResponse[T any](query string, cache pokecache.Cache) (T, error) {
	t := *new(T)

	body, err := getPokeapiJSONResponse(query, cache)
	if err != nil {
		return t, err
	}

	err = json.Unmarshal(body, &t)
	return t, err
}

var pokeapiLocationURL string = "https://pokeapi.co/api/v2/location-area/"

func locationURL(query string) string {
	return fmt.Sprintf("%s%s", pokeapiLocationURL, query)
}

func getPokeapiJSONResponse(query string, cache pokecache.Cache) ([]byte, error) {
	body, exists := cache.Get(query)
	if exists {
		return body, nil
	}
	resp, err := http.Get(query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	cache.Add(query, respBody)
	return respBody, nil
}

type LocationList struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type pokeapiResponse interface {
	LocationList | PokeapiLocation
}

// Generated with https://mholt.github.io/json-to-go/ from the example at https://pokeapi.co/docs/v2#location-areas
type PokeapiLocation struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int   `json:"min_level"`
				MaxLevel        int   `json:"max_level"`
				ConditionValues []any `json:"condition_values"`
				Chance          int   `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}
