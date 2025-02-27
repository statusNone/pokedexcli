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
			input:    " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  too much space  ",
			expected: []string{"too", "much", "space"},
		},
		{
			input:    " ",
			expected: []string{},
		},
		{
			input:    "GOlAnG Is AwEsOmE",
			expected: []string{"golang", "is", "awesome"},
		},
		{
			input:    "!@#$%^&*() hello!! world??",
			expected: []string{"!@#$%^&*()", "hello!!", "world??"},
		},
		{
			input:    "こんにちは 世界",
			expected: []string{"こんにちは", "世界"},
		},
		{
			input:    "hello\tworld\nnew\tline",
			expected: []string{"hello", "world", "new", "line"},
		},
		{
			input:    "!!!hello... world###",
			expected: []string{"!!!hello...", "world###"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("failed: input=%q, expected length %d, got length %d", c.input, len(c.expected), len(actual))
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("failed: input=%q at index %d, expected %q, got %q", c.input, i, expectedWord, word)
			}
		}
	}
}
