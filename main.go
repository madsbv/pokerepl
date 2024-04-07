package main

import (
	"bufio"
	"fmt"
	"github.com/madsbv/pokerepl/internal/pokeapi"
	"github.com/madsbv/pokerepl/internal/pokecache"
	"os"
	"time"
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
		commands(input).callback(&config)
		if !config.running {
			break
		}
	}
}

type command struct {
	name        string
	description string
	callback    func(*config)
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

func commandMap(c *config) {
	mapLocationsPage(c, c.next)
}

func commandMapb(c *config) {
	mapLocationsPage(c, c.prev)
}

func mapLocationsPage(c *config, dest *string) {
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

func commandHelp(c *config) {
	fmt.Println("Help menu")
}

func commandExit(c *config) {
	c.running = false
}
