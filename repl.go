package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/AlexSkr96/PokedexCLI/internal/pokecache"
)

type CLICommand struct {
	Name        string
	Description string
	Callback    func(*Config, []string) error
}

type Config struct {
	Previous string
	Next     string
}

type Pokemon struct {
	Height         int
	Weight         int
	HP             int
	Attack         int
	Defense        int
	SpecialAttack  int
	SpecialDefence int
	Speed          int
	Types          []string
}

var cache = pokecache.NewCache(5 * time.Second)
var pokedex = map[string]Pokemon{}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func commandExit(config *Config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config, args []string) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, cliCommand := range cliCommands {
		fmt.Printf("%v: %v\n", cliCommand.Name, cliCommand.Description)
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

func getMap(config *Config, dir bool) error {
	var data map[string]any
	var newNext string
	var newPrevious string
	var err error

	if dir {
		newPrevious = config.Next
		data, err = getFromAPI(newPrevious)
	} else {
		newNext = config.Previous
		data, err = getFromAPI(newNext)
	}
	if err != nil {
		return err
	}

	if dir {
		if data["next"] != nil {
			newNext = data["next"].(string)
		} else {
			newNext = ""
		}
	} else {
		if data["previous"] != nil {
			newPrevious = data["previous"].(string)
		} else {
			newPrevious = ""
		}
	}

	config.Next = newNext
	config.Previous = newPrevious

	for _, location := range data["results"].([]any) {
		convertedLocation := location.(map[string]any)
		fmt.Printf("  -%v\n", convertedLocation["name"])
	}

	return nil
}

func commandMap(config *Config, args []string) error {
	if config.Next == "" {
		fmt.Print("you're on the last page")
		return nil
	}
	return getMap(config, true)
}

func commandBMap(config *Config, args []string) error {
	if config.Previous == "" {
		fmt.Print("you're on the first page")
		return nil
	}
	return getMap(config, false)
}

func commandExplore(config *Config, args []string) error {
	fmt.Printf("Exploring %v...\n", args[0])
	fullUrl := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%v/", args[0])
	data, err := getFromAPI(fullUrl)
	if err != nil {
		return err
	}

	fmt.Printf("Found Pokemon:\n")
	for _, encounter := range data["pokemon_encounters"].([]any) {
		pokemon := encounter.(map[string]any)["pokemon"]
		name := pokemon.(map[string]any)["name"]
		fmt.Printf("  -%v\n", name)
	}

	return nil
}

func commandCatch(config *Config, args []string) error {
	fmt.Printf("Throwing a Pokeball at %v...\n", args[0])
	fullUrl := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%v/", args[0])
	data, err := getFromAPI(fullUrl)
	if err != nil {
		return err
	}

	if rand.Int()%150 >= int(data["base_experience"].(float64)) {
		pokemon := Pokemon{
			Height: int(data["height"].(float64)),
			Weight: int(data["weight"].(float64)),
			Types:  []string{},
		}
		for _, stat := range data["stats"].([]any) {
			mapStat := stat.(map[string]any)
			subStat := mapStat["stat"].(map[string]any)
			switch subStat["name"].(string) {
			case "hp":
				pokemon.HP = int(mapStat["base_stat"].(float64))
			case "attack":
				pokemon.Attack = int(mapStat["base_stat"].(float64))
			case "defense":
				pokemon.Defense = int(mapStat["base_stat"].(float64))
			case "special-attack":
				pokemon.SpecialAttack = int(mapStat["base_stat"].(float64))
			case "special-defense":
				pokemon.SpecialDefence = int(mapStat["base_stat"].(float64))
			case "speed":
				pokemon.Speed = int(mapStat["base_stat"].(float64))
			}
		}
		for _, typ := range data["types"].([]any) {
			typeMap := typ.(map[string]any)
			typeName := typeMap["type"].(map[string]any)["name"].(string)
			pokemon.Types = append(pokemon.Types, typeName)
		}
		pokedex[args[0]] = pokemon
		fmt.Printf("%v was caught!\n", args[0])
	} else {
		fmt.Printf("%v escaped!\n", args[0])
	}
	return nil
}
