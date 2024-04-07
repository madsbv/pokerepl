package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var pokeapiLocationURL string = "https://pokeapi.co/api/v2/location-area/"

func locationURL(query string) string {
	return fmt.Sprintf("%s%s", pokeapiLocationURL, query)
}

func GetLocations(q *string) (LocationList, error) {
	locations := LocationList{}

	query := ""
	// TODO: Do we want error handling instead of default behaviour here? I worry that default behaviour will break the expectation that prev and next are reversible, e.g., if we get to the end of the list and then `next` takes us back to the beginning.
	if q == nil {
		query = pokeapiLocationURL
	} else {
		query = *q
	}
	resp, err := http.Get(query)
	if err != nil {
		return locations, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return locations, err
	}

	err = json.Unmarshal(body, &locations)
	if err != nil {
		return locations, err
	}

	return locations, nil

	// locationNames := make([]string, 20)
	// for _, location := range locations.Results {
	// 	locationNames = append(locationNames, location.Name)
	// }
	// return locationNames, nil
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

// Generated with https://mholt.github.io/json-to-go/ from the example at https://pokeapi.co/docs/v2#location-areas
type pokeapiLocation struct {
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
