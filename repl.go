package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type CLICommand struct {
	name        string
	description string
	callback    func(*Config) error
}

type Config struct {
	previous string
	next     string
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func commandExit(config *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, cliCommand := range cliCommands {
		fmt.Printf("%v: %v\n", cliCommand.name, cliCommand.description)
	}
	return nil
}

func getFromAPI(url string) (map[string]any, error) {
	res, err := http.Get(url)
	if err != nil {
		return map[string]any{}, fmt.Errorf("error while getting '%v': %v", url, err)
	}
	defer res.Body.Close()

	var data map[string]any
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return map[string]any{}, fmt.Errorf("error while decoding JSON response: %v", err)
	}

	return data, nil
}

func getMap(config *Config, dir bool) (err error) {
	var data map[string]any
	if dir {
		data, err = getFromAPI(config.next)
	} else {
		data, err = getFromAPI(config.previous)
	}
	if err != nil {
		return err
	}

	if data["next"] != nil {
		config.next = data["next"].(string)
	} else {
		config.next = ""
	}
	if data["previous"] != nil {
		config.previous = data["previous"].(string)
	} else {
		config.previous = ""
	}

	for _, location := range data["results"].([]any) {
		convertedLocation := location.(map[string]any)
		fmt.Printf("%v\n", convertedLocation["name"])
	}

	return nil
}

func commandMap(config *Config) error {
	return getMap(config, true)
}

func commandBMap(config *Config) error {
	return getMap(config, false)
}
