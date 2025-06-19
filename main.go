package main

import (
	"bufio"
	"fmt"
	"os"
)

var cliCommands = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
}

func main() {
	cliCommands["help"] = cliCommand{
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
			err := cliCommand.callback()
			if err != nil {
				fmt.Printf("ERROR: %v", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
