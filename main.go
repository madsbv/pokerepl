package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/madsbv/pokerepl/internal/pokeapi"
	"github.com/madsbv/pokerepl/internal/pokecache"
)

func main() {
	prompt := "pokedex > "
	scanner := bufio.NewScanner(os.Stdin)
	cacheInterval := 60 * 5 * time.Second
	config := config{nil, nil, true, pokecache.New(cacheInterval), make(map[string]pokeapi.PokeapiPokemon)}
	for {
		fmt.Print(prompt)
		scanner.Scan()
		input := scanner.Text()
		args := strings.Split(input, " ")
		cmd := args[0]
		commands(cmd).callback(&config, args[1:])
		if !config.running {
			break
		}
	}
}

type command struct {
	name        string
	description string
	callback    func(*config, []string)
}

func commands(input string) command {
	commands := map[string]command{
		"help": {
			name:        "help",
			description: "Displays this help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Go forwards and display map",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Go back and display map",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explore an area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a Pokemon!",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a Pokemon you have caught",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List the Pokemon you have caught",
			callback:    commandPokedex,
		},
	}

	// Command aliases
	if input == "q" {
		input = "exit"
	}

	com, ok := commands[input]
	if ok {
		return com
	} else {
		fmt.Println("Unknown command")
		return commands["help"]
	}
}

type config struct {
	next    *string
	prev    *string
	running bool
	cache   pokecache.Cache
	pokeman map[string]pokeapi.PokeapiPokemon
}

func commandHelp(c *config, _ []string) {
	fmt.Println("Help menu")
	// TODO: Expand on this
}

func commandExit(c *config, _ []string) {
	c.running = false
}

func commandMap(c *config, _ []string) {
	printLocationsPage(c, c.next)
}

func commandMapb(c *config, _ []string) {
	printLocationsPage(c, c.prev)
}

func printLocationsPage(c *config, dest *string) {
	p, err := pokeapi.GetLocations(dest, c.cache)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	locationNames := make([]string, 0, 20)
	for _, location := range p.Results {
		locationNames = append(locationNames, location.Name)
	}
	c.next = p.Next
	c.prev = p.Previous
	for _, n := range locationNames {
		fmt.Println(n)
	}
}

func commandExplore(c *config, args []string) {
	if len(args) == 0 {
		fmt.Println("Enter the name of an area to explore it further!")
		return
	}

	name := args[0]

	areaDetails, err := pokeapi.GetLocationDetails(name, c.cache)
	if err != nil {
		fmt.Printf("Something went wrong while exploring %v\n", name)
	}

	fmt.Printf("Exploring %v...\n", name)
	fmt.Println("Found Pokemon:")
	for _, encounter := range areaDetails.PokemonEncounters {
		fmt.Printf(" - %v\n", encounter.Pokemon.Name)
	}
}

func commandCatch(c *config, args []string) {
	if len(args) == 0 {
		fmt.Println("Enter the name of a Pokemon to try to catch it!")
		return
	}
	pokemon, err := pokeapi.GetPokemonDetails(args[0], c.cache)
	if err != nil {
		fmt.Printf("Something went wrong while trying to catch %v.\n", args[0])
		return
	}

	name, exp := pokemon.Name, pokemon.BaseExperience
	// TODO: Mewtwo has base experience 340, change the rng?
	if exp > 200 {
		fmt.Printf("%v has base experience %v\n", name, exp)
	}

	roll := rand.Intn(201)
	if roll >= exp {
		fmt.Printf("You caught a %v!\n", name)
		c.pokeman[name] = pokemon
	} else {
		fmt.Printf("%v got away!\n", name)
	}
}

func commandInspect(c *config, args []string) {
	if len(args) == 0 {
		fmt.Println("Enter the name of a Pokemon to try to inspect")
		return
	}
	name := args[0]
	pokemon, exists := c.pokeman[name]
	if !exists {
		fmt.Printf("You have not caught a %v\n", name)
	}
	fmt.Printf("Name: %v\n", pokemon.Name)
	fmt.Printf("Height: %v\n", pokemon.Height)
	fmt.Printf("Weight: %v\n", pokemon.Weight)
	fmt.Printf("Stats:\n")
	for _, s := range pokemon.Stats {
		fmt.Printf("  - %v: %v\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Printf("Types:\n")
	for _, t := range pokemon.Types {
		fmt.Printf("  - %v\n", t.Type.Name)
	}
}

func commandPokedex(c *config, args []string) {
	fmt.Println("Your Pokedex:")
	for k, _ := range c.pokeman {
		fmt.Printf("  - %v\n", k)
	}
}
