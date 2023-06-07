package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextToLines(t *testing.T) {
	for _, tt := range []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "10 chars one word",
			input:    "0123456789",
			expected: []string{"0123456789"},
		},
		/*
			{
				name:     "11 chars one word",
				input:    "01234567891",
				expected: []string{"0123456789", "1"},
			},
		*/
		{
			name:     "10 chars two words",
			input:    "012345678 9",
			expected: []string{"012345678", "9"},
		},
		{
			name:     "two words with more than 10 characters",
			input:    "Hello world!",
			expected: []string{"Hello", "world!"},
		},
		{
			name:     "with linebreak",
			input:    "a\na",
			expected: []string{"a", "", "a"},
		},
		{
			name:     "with multiple linebreaks",
			input:    "a\n\n\na",
			expected: []string{"a", "", "a"},
		},
		/*
			{
				name:     "with extra linebreaks",
				input:    "\na\na",
				expected: []string{"a"},
			},
		*/
		/*
			{
				name:     "with extra spaces",
				input:    " a ",
				expected: []string{"a"},
			},
		*/
		{
			name:     "11 chars one word, with trailing space",
			input:    "0123456789 ",
			expected: []string{"0123456789"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			actual := TextToLines(tt.input, 10)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
