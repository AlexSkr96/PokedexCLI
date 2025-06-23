package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/AlexSkr96/PokedexCLI/internal/pokecache"
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

var cache = pokecache.NewCache(5 * time.Second)

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
	jsonData, exists := cache.Get(url)
	if !exists {
		res, err := http.Get(url)
		if err != nil {
			return map[string]any{}, fmt.Errorf("error getting \"%v\": %v", url, err)
		}
		defer res.Body.Close()

		jsonData, err = io.ReadAll(res.Body)
		if err != nil {
			return map[string]any{}, fmt.Errorf("error reading response body from \"%v\": %v", url, err)
		}

		cache.Add(url, jsonData)
	}

	var data map[string]any
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		return map[string]any{}, fmt.Errorf("error unmarshalling data from \"%v\": %v", url, err)
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
