package pokeapi

import (
	"encoding/json"
	"github.com/madsbv/pokerepl/internal/pokecache"
	"io"
	"net/http"
)

func getParsedResponse[T any](query string, cache pokecache.Cache) (T, error) {
	t := *new(T)

	body, err := getPokeapiJSONResponse(query, cache)
	if err != nil {
		return t, err
	}

	err = json.Unmarshal(body, &t)
	return t, err
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

type PokeapiResponse interface {
	LocationList | PokeapiLocation
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

type PokeapiPokemon struct {
	ID             int    `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	BaseExperience int    `json:"base_experience,omitempty"`
	Height         int    `json:"height,omitempty"`
	IsDefault      bool   `json:"is_default,omitempty"`
	Order          int    `json:"order,omitempty"`
	Weight         int    `json:"weight,omitempty"`
	Abilities      []struct {
		IsHidden bool `json:"is_hidden,omitempty"`
		Slot     int  `json:"slot,omitempty"`
		Ability  struct {
			Name string `json:"name,omitempty"`
			URL  string `json:"url,omitempty"`
		} `json:"ability,omitempty"`
	} `json:"abilities,omitempty"`
	Forms []struct {
		Name string `json:"name,omitempty"`
		URL  string `json:"url,omitempty"`
	} `json:"forms,omitempty"`
	GameIndices []struct {
		GameIndex int `json:"game_index,omitempty"`
		Version   struct {
			Name string `json:"name,omitempty"`
			URL  string `json:"url,omitempty"`
		} `json:"version,omitempty"`
	} `json:"game_indices,omitempty"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name,omitempty"`
			URL  string `json:"url,omitempty"`
		} `json:"item,omitempty"`
		VersionDetails []struct {
			Rarity  int `json:"rarity,omitempty"`
			Version struct {
				Name string `json:"name,omitempty"`
				URL  string `json:"url,omitempty"`
			} `json:"version,omitempty"`
		} `json:"version_details,omitempty"`
	} `json:"held_items,omitempty"`
	LocationAreaEncounters string `json:"location_area_encounters,omitempty"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name,omitempty"`
			URL  string `json:"url,omitempty"`
		} `json:"move,omitempty"`
		VersionGroupDetails []struct {
			LevelLearnedAt int `json:"level_learned_at,omitempty"`
			VersionGroup   struct {
				Name string `json:"name,omitempty"`
				URL  string `json:"url,omitempty"`
			} `json:"version_group,omitempty"`
			MoveLearnMethod struct {
				Name string `json:"name,omitempty"`
				URL  string `json:"url,omitempty"`
			} `json:"move_learn_method,omitempty"`
		} `json:"version_group_details,omitempty"`
	} `json:"moves,omitempty"`
	Species struct {
		Name string `json:"name,omitempty"`
		URL  string `json:"url,omitempty"`
	} `json:"species,omitempty"`
	Sprites struct {
		BackDefault      string `json:"back_default,omitempty"`
		BackFemale       any    `json:"back_female,omitempty"`
		BackShiny        string `json:"back_shiny,omitempty"`
		BackShinyFemale  any    `json:"back_shiny_female,omitempty"`
		FrontDefault     string `json:"front_default,omitempty"`
		FrontFemale      any    `json:"front_female,omitempty"`
		FrontShiny       string `json:"front_shiny,omitempty"`
		FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
		Other            struct {
			DreamWorld struct {
				FrontDefault string `json:"front_default,omitempty"`
				FrontFemale  any    `json:"front_female,omitempty"`
			} `json:"dream_world,omitempty"`
			Home struct {
				FrontDefault     string `json:"front_default,omitempty"`
				FrontFemale      any    `json:"front_female,omitempty"`
				FrontShiny       string `json:"front_shiny,omitempty"`
				FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
			} `json:"home,omitempty"`
			OfficialArtwork struct {
				FrontDefault string `json:"front_default,omitempty"`
				FrontShiny   string `json:"front_shiny,omitempty"`
			} `json:"official-artwork,omitempty"`
			Showdown struct {
				BackDefault      string `json:"back_default,omitempty"`
				BackFemale       any    `json:"back_female,omitempty"`
				BackShiny        string `json:"back_shiny,omitempty"`
				BackShinyFemale  any    `json:"back_shiny_female,omitempty"`
				FrontDefault     string `json:"front_default,omitempty"`
				FrontFemale      any    `json:"front_female,omitempty"`
				FrontShiny       string `json:"front_shiny,omitempty"`
				FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
			} `json:"showdown,omitempty"`
		} `json:"other,omitempty"`
		Versions struct {
			GenerationI struct {
				RedBlue struct {
					BackDefault  string `json:"back_default,omitempty"`
					BackGray     string `json:"back_gray,omitempty"`
					FrontDefault string `json:"front_default,omitempty"`
					FrontGray    string `json:"front_gray,omitempty"`
				} `json:"red-blue,omitempty"`
				Yellow struct {
					BackDefault  string `json:"back_default,omitempty"`
					BackGray     string `json:"back_gray,omitempty"`
					FrontDefault string `json:"front_default,omitempty"`
					FrontGray    string `json:"front_gray,omitempty"`
				} `json:"yellow,omitempty"`
			} `json:"generation-i,omitempty"`
			GenerationIi struct {
				Crystal struct {
					BackDefault  string `json:"back_default,omitempty"`
					BackShiny    string `json:"back_shiny,omitempty"`
					FrontDefault string `json:"front_default,omitempty"`
					FrontShiny   string `json:"front_shiny,omitempty"`
				} `json:"crystal,omitempty"`
				Gold struct {
					BackDefault  string `json:"back_default,omitempty"`
					BackShiny    string `json:"back_shiny,omitempty"`
					FrontDefault string `json:"front_default,omitempty"`
					FrontShiny   string `json:"front_shiny,omitempty"`
				} `json:"gold,omitempty"`
				Silver struct {
					BackDefault  string `json:"back_default,omitempty"`
					BackShiny    string `json:"back_shiny,omitempty"`
					FrontDefault string `json:"front_default,omitempty"`
					FrontShiny   string `json:"front_shiny,omitempty"`
				} `json:"silver,omitempty"`
			} `json:"generation-ii,omitempty"`
			GenerationIii struct {
				Emerald struct {
					FrontDefault string `json:"front_default,omitempty"`
					FrontShiny   string `json:"front_shiny,omitempty"`
				} `json:"emerald,omitempty"`
				FireredLeafgreen struct {
					BackDefault  string `json:"back_default,omitempty"`
					BackShiny    string `json:"back_shiny,omitempty"`
					FrontDefault string `json:"front_default,omitempty"`
					FrontShiny   string `json:"front_shiny,omitempty"`
				} `json:"firered-leafgreen,omitempty"`
				RubySapphire struct {
					BackDefault  string `json:"back_default,omitempty"`
					BackShiny    string `json:"back_shiny,omitempty"`
					FrontDefault string `json:"front_default,omitempty"`
					FrontShiny   string `json:"front_shiny,omitempty"`
				} `json:"ruby-sapphire,omitempty"`
			} `json:"generation-iii,omitempty"`
			GenerationIv struct {
				DiamondPearl struct {
					BackDefault      string `json:"back_default,omitempty"`
					BackFemale       any    `json:"back_female,omitempty"`
					BackShiny        string `json:"back_shiny,omitempty"`
					BackShinyFemale  any    `json:"back_shiny_female,omitempty"`
					FrontDefault     string `json:"front_default,omitempty"`
					FrontFemale      any    `json:"front_female,omitempty"`
					FrontShiny       string `json:"front_shiny,omitempty"`
					FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
				} `json:"diamond-pearl,omitempty"`
				HeartgoldSoulsilver struct {
					BackDefault      string `json:"back_default,omitempty"`
					BackFemale       any    `json:"back_female,omitempty"`
					BackShiny        string `json:"back_shiny,omitempty"`
					BackShinyFemale  any    `json:"back_shiny_female,omitempty"`
					FrontDefault     string `json:"front_default,omitempty"`
					FrontFemale      any    `json:"front_female,omitempty"`
					FrontShiny       string `json:"front_shiny,omitempty"`
					FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
				} `json:"heartgold-soulsilver,omitempty"`
				Platinum struct {
					BackDefault      string `json:"back_default,omitempty"`
					BackFemale       any    `json:"back_female,omitempty"`
					BackShiny        string `json:"back_shiny,omitempty"`
					BackShinyFemale  any    `json:"back_shiny_female,omitempty"`
					FrontDefault     string `json:"front_default,omitempty"`
					FrontFemale      any    `json:"front_female,omitempty"`
					FrontShiny       string `json:"front_shiny,omitempty"`
					FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
				} `json:"platinum,omitempty"`
			} `json:"generation-iv,omitempty"`
			GenerationV struct {
				BlackWhite struct {
					Animated struct {
						BackDefault      string `json:"back_default,omitempty"`
						BackFemale       any    `json:"back_female,omitempty"`
						BackShiny        string `json:"back_shiny,omitempty"`
						BackShinyFemale  any    `json:"back_shiny_female,omitempty"`
						FrontDefault     string `json:"front_default,omitempty"`
						FrontFemale      any    `json:"front_female,omitempty"`
						FrontShiny       string `json:"front_shiny,omitempty"`
						FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
					} `json:"animated,omitempty"`
					BackDefault      string `json:"back_default,omitempty"`
					BackFemale       any    `json:"back_female,omitempty"`
					BackShiny        string `json:"back_shiny,omitempty"`
					BackShinyFemale  any    `json:"back_shiny_female,omitempty"`
					FrontDefault     string `json:"front_default,omitempty"`
					FrontFemale      any    `json:"front_female,omitempty"`
					FrontShiny       string `json:"front_shiny,omitempty"`
					FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
				} `json:"black-white,omitempty"`
			} `json:"generation-v,omitempty"`
			GenerationVi struct {
				OmegarubyAlphasapphire struct {
					FrontDefault     string `json:"front_default,omitempty"`
					FrontFemale      any    `json:"front_female,omitempty"`
					FrontShiny       string `json:"front_shiny,omitempty"`
					FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
				} `json:"omegaruby-alphasapphire,omitempty"`
				XY struct {
					FrontDefault     string `json:"front_default,omitempty"`
					FrontFemale      any    `json:"front_female,omitempty"`
					FrontShiny       string `json:"front_shiny,omitempty"`
					FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
				} `json:"x-y,omitempty"`
			} `json:"generation-vi,omitempty"`
			GenerationVii struct {
				Icons struct {
					FrontDefault string `json:"front_default,omitempty"`
					FrontFemale  any    `json:"front_female,omitempty"`
				} `json:"icons,omitempty"`
				UltraSunUltraMoon struct {
					FrontDefault     string `json:"front_default,omitempty"`
					FrontFemale      any    `json:"front_female,omitempty"`
					FrontShiny       string `json:"front_shiny,omitempty"`
					FrontShinyFemale any    `json:"front_shiny_female,omitempty"`
				} `json:"ultra-sun-ultra-moon,omitempty"`
			} `json:"generation-vii,omitempty"`
			GenerationViii struct {
				Icons struct {
					FrontDefault string `json:"front_default,omitempty"`
					FrontFemale  any    `json:"front_female,omitempty"`
				} `json:"icons,omitempty"`
			} `json:"generation-viii,omitempty"`
		} `json:"versions,omitempty"`
	} `json:"sprites,omitempty"`
	Cries struct {
		Latest string `json:"latest,omitempty"`
		Legacy string `json:"legacy,omitempty"`
	} `json:"cries,omitempty"`
	Stats []struct {
		BaseStat int `json:"base_stat,omitempty"`
		Effort   int `json:"effort,omitempty"`
		Stat     struct {
			Name string `json:"name,omitempty"`
			URL  string `json:"url,omitempty"`
		} `json:"stat,omitempty"`
	} `json:"stats,omitempty"`
	Types []struct {
		Slot int `json:"slot,omitempty"`
		Type struct {
			Name string `json:"name,omitempty"`
			URL  string `json:"url,omitempty"`
		} `json:"type,omitempty"`
	} `json:"types,omitempty"`
	PastTypes []struct {
		Generation struct {
			Name string `json:"name,omitempty"`
			URL  string `json:"url,omitempty"`
		} `json:"generation,omitempty"`
		Types []struct {
			Slot int `json:"slot,omitempty"`
			Type struct {
				Name string `json:"name,omitempty"`
				URL  string `json:"url,omitempty"`
			} `json:"type,omitempty"`
		} `json:"types,omitempty"`
	} `json:"past_types,omitempty"`
}
