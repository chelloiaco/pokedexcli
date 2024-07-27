package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type cliCommand struct {
	name        string
	description string
	Callback    func(string) error
}

/*
List of all commands availabe in the program.
Contains their name, description and Callback function
*/
func GetCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays this message",
			Callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Lists the 20 areas forward available to 'explore'",
			Callback:    commandMapForward,
		},
		"mapb": {
			name:        "mapb",
			description: "Lists the 20 areas backward available to 'explore'",
			Callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "Lists all of the Pokemon of a given area",
			Callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Tries to catch a Pokemon!",
			Callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspects a given pokemon if present in the pokedex",
			Callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Prints out the Pokemons you have caught so far",
			Callback:    commandPokedex,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			Callback:    commandExit,
		},
	}
}

/* Callback functions below */

/*
commandHelp
Lists all comands available in the program for ease of use.
*/
func commandHelp(_ string) error {
	fmt.Println(`
Welcome to the Pokedex!
Usage:`)
	fmt.Println("")
	for _, v := range GetCommands() {
		fmt.Printf("%s: %s\n", v.name, v.description)
	}
	fmt.Println("")
	return nil
}

/*
commandExit
Exits the program.
*/
func commandExit(_ string) error {
	os.Exit(0)
	return nil
}

/*
commandMapForward
Lists to the user the 20 next location areas to explore.
*/
func commandMapForward(_ string) error {
	if nextUrl == "" {
		fmt.Println("Error: cannot map futher")
		return nil
	}
	printLocationAreas(nextUrl)

	return nil
}

/*
commandMapBack
Lists to the user the 20 previous location areas to explore.
*/
func commandMapBack(_ string) error {
	if previousUrl == "" {
		fmt.Println("Error: cannot map back")
		return nil
	}
	printLocationAreas(previousUrl)

	return nil
}

/*
commandExplore
Explores the given area, listing all pokemons that are available for catching
*/
func commandExplore(name string) error {
	if name == "" {
		fmt.Println("Please provide the name of the area to explore.")
		return nil
	}
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%v", name)

	// Fetch data from either cache or GET
	data, err := FetchData(url)
	if err != nil {
		return nil
	}

	var locationarea LocationArea
	err = json.Unmarshal(*data, &locationarea)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Exploring %s...\n", locationarea.Name)
	fmt.Println("Found Pokemon:")
	for i := range locationarea.PokemonEncounters {
		pokemonEncounter := locationarea.PokemonEncounters[i]
		fmt.Printf(" - %s\n", pokemonEncounter.Pokemon.Name)
	}

	return nil
}

/*
commandCatch
Tries to catch the given Pokemon name. If succeedes, adds that to the pokemon
map
*/
func commandCatch(name string) error {
	if name == "" {
		fmt.Println("Please provide the name of the Pokemon to try to catch.")
		return nil
	}
	_, ok := pokedex[name]
	if ok {
		fmt.Printf("%v is already on the Pokedex!\n", name)
		return nil
	}
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%v", name)

	// Fetch data form either cache or GET
	data, err := FetchData(url)
	if err != nil {
		return nil
	}

	var pokemon Pokemon
	err = json.Unmarshal(*data, &pokemon)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Throwing a Pokeball at %v...\n", pokemon.Name)

	// Try to catch Pokemon
	success := catchAttempt(pokemon)
	if success {
		fmt.Printf("%v was caught!\n", pokemon.Name)
		fmt.Println("You may now inspect it with the 'inspect' command.")
		// Add it to the Pokedex
		pokedex[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%v escaped!\n", pokemon.Name)
	}

	return nil
}

/*
commandInspect
Inspects the given pokemon name if it exists in the pokedex.
*/
func commandInspect(name string) error {
	if name == "" {
		fmt.Println("Please provide the name of the Pokemon to inspect.")
		return nil
	}
	pokemon, ok := pokedex[name]
	if !ok {
		fmt.Printf("%v was not caught yet!\n", name)
		return nil
	}

	// Print basic info
	fmt.Printf("Name: %v\n", pokemon.Name)
	fmt.Printf("Height: %v\n", pokemon.Height)
	fmt.Printf("Weight: %v\n", pokemon.Weight)

	// Print Stats
	fmt.Println("Stats:")
	for i := range pokemon.Stats {
		baseStat := pokemon.Stats[i].BaseStat
		statName := pokemon.Stats[i].Stat.Name
		fmt.Printf("  - %v: %v\n", baseStat, statName)
	}

	// Print Types
	fmt.Println("Types:")
	for i := range pokemon.Types {
		fmt.Printf("  - %v\n", pokemon.Types[i].Type.Name)
	}

	return nil
}

/*
commandPokedex
Inspects the given pokemon name if it exists in the pokedex.
*/
func commandPokedex(_ string) error {
	if len(pokedex) < 1 {
		fmt.Println("You haven't caught any Pokemon yet!")
		return nil
	}

	for i := range pokedex {
		fmt.Printf("  - %v\n", pokedex[i].Name)
	}

	return nil
}
