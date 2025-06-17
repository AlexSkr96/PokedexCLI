package main

import (
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
