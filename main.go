package main

// test commit
import (
	"bufio"
	"fmt"
	"os"
)

var cliCommands = map[string]CLICommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
	"map": {
		name:        "map",
		description: "Get next 20 locations of Pokemon world",
		callback:    commandMap,
	},
	"bmap": {
		name:        "bmap",
		description: "Get previous 20 locations of Pokemon world",
		callback:    commandBMap,
	},
}

var locationConfig = Config{
	next:     "https://pokeapi.co/api/v2/location-area/",
	previous: "",
}

func main() {
	// add help here to avoid error of cyclycal use (commandHelp func <> cliCommand struct)
	cliCommands["help"] = CLICommand{
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
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
			err := cliCommand.callback(&locationConfig)
			if err != nil {
				fmt.Printf("ERROR: %v", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
