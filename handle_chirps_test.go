package main

import (
	"testing"
)

func TestCleanMsg(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "Hello, Word!",
			expected: "Hello, Word!",
		},
		{
			input:    "You sharbert!",
			expected: "You sharbert!",
		},
		{
			input:    "Never Fornax away",
			expected: "Never **** away",
		},
	}
	for _, c := range cases {
		actual := cleanMsg(c.input)
		if actual != c.expected {
			t.Errorf(
				"Incorect cleaning for: %s\n\tgot: %s\n\texpected: %s\n",
				c.input,
				actual,
				c.expected,
			)
		}
	}
}

func TestCheckProfanity(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "cat",
			expected: "cat",
		},
		{
			input:    "kerfuffle",
			expected: "****",
		},
		{
			input:    "kerfufle",
			expected: "kerfufle",
		},
		{
			input:    "sharbert",
			expected: "****",
		},
		{
			input:    "fornax",
			expected: "****",
		},
	}
	for _, c := range cases {
		actual := checkProfanity(c.input)
		if actual != c.expected {
			t.Errorf(
				"Incorrect profanity check for: %s\n\tgot: %s\n\texpected: %s\n",
				c.input,
				actual,
				c.expected,
			)
		}
	}
}
