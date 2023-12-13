package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextToLines(t *testing.T) {
	width := 10

	t.Run("shorter than width should not split", func(t *testing.T) {
		i := "test"
		e := []string{"test"}
		assert.Equal(t, e, TextToLines(i, width))
	})

	t.Run("should split cleanly on space", func(t *testing.T) {
		i := "abcde fghij" // 11 chars
		e := []string{"abcde", "fghij"}
		assert.Equal(t, e, TextToLines(i, width))
	})

	t.Run("space should not split if not long enough", func(t *testing.T) {
		i := "a b c"
		e := []string{"a b c"}
		assert.Equal(t, e, TextToLines(i, width))
	})

	t.Run("should split hard once width has been reached", func(t *testing.T) {
		t.Skip()             // TODO
		i := "abcdefghijklm" // 13 chars
		e := []string{"abcdefghij", "klm"}
		assert.Equal(t, e, TextToLines(i, width))
	})

	t.Run("linebreaks should add empty lines", func(t *testing.T) {
		i := "a\nb"
		e := []string{"a", "", "b"}
		assert.Equal(t, e, TextToLines(i, width))
	})

	t.Run("multiple linebreaks should not add multiple empty lines", func(t *testing.T) {
		i := "a\n\nb"
		e := []string{"a", "", "b"}
		assert.Equal(t, e, TextToLines(i, width))
	})

	t.Run("leading linebreak should be trimmed", func(t *testing.T) {
		i := "\na"
		e := []string{"a"}
		assert.Equal(t, e, TextToLines(i, width))
	})

	t.Run("trailing linebreak should be trimmed", func(t *testing.T) {
		i := "a\n"
		e := []string{"a"}
		assert.Equal(t, e, TextToLines(i, width))
	})

	t.Run("lorem ipsum example", func(t *testing.T) {
		i := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
		e := []string{
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed",
			"do eiusmod tempor incididunt ut labore et dolore magna",
			"aliqua. Ut enim ad minim veniam, quis nostrud exercitation",
			"ullamco laboris nisi ut aliquip ex ea commodo consequat.",
			//"", // TODO: Implement nice paragraph splitting
			"Duis aute irure dolor in reprehenderit in voluptate velit",
			"esse cillum dolore eu fugiat nulla pariatur. Excepteur sint",
			"occaecat cupidatat non proident, sunt in culpa qui officia",
			"deserunt mollit anim id est laborum.",
		}
		assert.Equal(t, e, TextToLines(i, 60))
	})
}
