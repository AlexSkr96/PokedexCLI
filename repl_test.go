package main

import (
	"fmt"
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Lorem Ipsum",
			expected: []string{"lorem", "ipsum"},
		},
		{
			input:    "Funny	little  bunnY",
			expected: []string{"funny", "little", "bunny"},
		},
		{
			input:    "This is the first line\n" + "And this is the second",
			expected: []string{"this", "is", "the", "first", "line", "and", "this", "is", "the", "second"},
		},
	}

	for _, cas := range cases {
		actual := cleanInput(cas.input)
		if len(actual) != len(cas.expected) {
			t.Errorf("expected %v word, got %v words", len(cas.expected), len(actual))
			t.Fail()
		}
		for i := range actual {
			if actual[i] != cas.expected[i] {
				t.Errorf("expected %v, got %v", cas.expected[i], actual[i])
				t.Fail()
			}
		}
	}
}

func HelpTestMap(testSubject, expectedResult map[string]any) error {
	for k, v := range testSubject {
		switch k {
		case "results":
			// Handle JSON unmarshaling where results is []interface{}
			actualResults := v.([]any)
			expectedResults := expectedResult[k].([]map[string]any)
			for i, el := range actualResults {
				actualMap := el.(map[string]any)
				err := HelpTestMap(actualMap, expectedResults[i])
				if err != nil {
					return err
				}
			}
		default:
			if v != expectedResult[k] {
				return fmt.Errorf("at \"%v\" expected %v, got %v", k, expectedResult[k], v)
			}
		}
	}
	return nil
}

func TestGetLocationsFromAPI(t *testing.T) {
	cases := []struct {
		input    string
		expected map[string]any
	}{
		{
			input: "https://pokeapi.co/api/v2/location-area/",
			expected: map[string]any{
				"count":    float64(1089),
				"next":     "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20",
				"previous": nil,
				"results": []map[string]any{
					{
						"name": "canalave-city-area",
						"url":  "https://pokeapi.co/api/v2/location-area/1/",
					},
					{
						"name": "eterna-city-area",
						"url":  "https://pokeapi.co/api/v2/location-area/2/",
					},
					{
						"name": "pastoria-city-area",
						"url":  "https://pokeapi.co/api/v2/location-area/3/",
					},
					{
						"name": "sunyshore-city-area",
						"url":  "https://pokeapi.co/api/v2/location-area/4/",
					},
					{
						"name": "sinnoh-pokemon-league-area",
						"url":  "https://pokeapi.co/api/v2/location-area/5/",
					},
					{
						"name": "oreburgh-mine-1f",
						"url":  "https://pokeapi.co/api/v2/location-area/6/",
					},
					{
						"name": "oreburgh-mine-b1f",
						"url":  "https://pokeapi.co/api/v2/location-area/7/",
					},
					{
						"name": "valley-windworks-area",
						"url":  "https://pokeapi.co/api/v2/location-area/8/",
					},
					{
						"name": "eterna-forest-area",
						"url":  "https://pokeapi.co/api/v2/location-area/9/",
					},
					{
						"name": "fuego-ironworks-area",
						"url":  "https://pokeapi.co/api/v2/location-area/10/",
					},
					{
						"name": "mt-coronet-1f-route-207",
						"url":  "https://pokeapi.co/api/v2/location-area/11/",
					},
					{
						"name": "mt-coronet-2f",
						"url":  "https://pokeapi.co/api/v2/location-area/12/",
					},
					{
						"name": "mt-coronet-3f",
						"url":  "https://pokeapi.co/api/v2/location-area/13/",
					},
					{
						"name": "mt-coronet-exterior-snowfall",
						"url":  "https://pokeapi.co/api/v2/location-area/14/",
					},
					{
						"name": "mt-coronet-exterior-blizzard",
						"url":  "https://pokeapi.co/api/v2/location-area/15/",
					},
					{
						"name": "mt-coronet-4f",
						"url":  "https://pokeapi.co/api/v2/location-area/16/",
					},
					{
						"name": "mt-coronet-4f-small-room",
						"url":  "https://pokeapi.co/api/v2/location-area/17/",
					},
					{
						"name": "mt-coronet-5f",
						"url":  "https://pokeapi.co/api/v2/location-area/18/",
					},
					{
						"name": "mt-coronet-6f",
						"url":  "https://pokeapi.co/api/v2/location-area/19/",
					},
					{
						"name": "mt-coronet-1f-from-exterior",
						"url":  "https://pokeapi.co/api/v2/location-area/20/",
					},
				},
			},
		},
	}

	for _, cas := range cases {
		actual, err := getFromAPI(cas.input)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		err = HelpTestMap(actual, cas.expected)
		if err != nil {
			t.Errorf("%v", err)
		}
	}
}
