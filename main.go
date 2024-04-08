package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/madsbv/pokerepl/internal/pokeapi"
	"github.com/madsbv/pokerepl/internal/pokecache"
)

func main() {
	prompt := "pokedex > "
	scanner := bufio.NewScanner(os.Stdin)
	cacheInterval := 5 * time.Second
	config := config{nil, nil, true, pokecache.New(cacheInterval)}
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
