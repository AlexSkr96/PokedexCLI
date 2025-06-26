package main

import (
	"bufio"
	"fmt"
	"os"
)

var cliCommands = map[string]CLICommand{
	"exit": {
		Name:        "exit",
		Description: "Exit the Pokedex",
		Callback:    commandExit,
	},
	"map": {
		Name:        "map",
		Description: "Get next 20 locations of Pokemon world",
		Callback:    commandMap,
	},
	"bmap": {
		Name:        "bmap",
		Description: "Get previous 20 locations of Pokemon world",
		Callback:    commandBMap,
	},
	"explore": {
		Name:        "explore",
		Description: "See a list of pokemon in selected area",
		Callback:    commandExplore,
	},
	"catch": {
		Name:        "catch",
		Description: "Attempt to catch a pokemon",
		Callback:    commandCatch,
	},
	"inspect": {
		Name:        "inspect",
		Description: "Look up info on caught pokemon",
		Callback:    commandInspect,
	},
	"pokedex": {
		Name:        "pokedex",
		Description: "Display a list of all caught pokemon",
		Callback:    commandPokedex,
	},
}

var locationConfig = Config{
	Next:     "https://pokeapi.co/api/v2/location-area/",
	Previous: "",
}

func main() {
	// add help here to avoid error of cyclycal use (commandHelp func <> cliCommand struct)
	cliCommands["help"] = CLICommand{
		Name:        "help",
		Description: "Displays a help message",
		Callback:    commandHelp,
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cleanedInput := cleanInput(input)
		var command string
		if len(cleanedInput) == 0 {
			fmt.Print("")
			continue
		} else {
			command = cleanedInput[0]
		}

		if cliCommand, exists := cliCommands[command]; exists {
			err := cliCommand.Callback(&locationConfig, cleanedInput[1:])
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
