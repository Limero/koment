package wordpress

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mustURL(s string) *url.URL {
	u, _ := url.Parse(s)
	return u
}

func TestToModelFlat(t *testing.T) {
	comments := []WordPressComment{
		{ID: 1, Parent: 0, AuthorName: "Alice", DateGmt: "2026-06-05T12:26:41", Content: struct {
			Rendered string `json:"rendered"`
		}{Rendered: "First"}},
		{ID: 2, Parent: 0, AuthorName: "Bob", DateGmt: "2026-06-05T12:27:41", Content: struct {
			Rendered string `json:"rendered"`
		}{Rendered: "Second"}},
	}

	posts := toModel(comments)
	assert.Len(t, posts, 2)
	assert.Equal(t, 0, posts[0].Depth)
	assert.Equal(t, 0, posts[1].Depth)
	assert.Equal(t, "Alice", posts[0].Author.Name)
	assert.Equal(t, "Bob", posts[1].Author.Name)
}

func TestToModelThreaded(t *testing.T) {
	comments := []WordPressComment{
		{ID: 1, Parent: 0, AuthorName: "Alice", DateGmt: "2026-06-05T12:26:41", Content: struct {
			Rendered string `json:"rendered"`
		}{Rendered: "Top"}},
		{ID: 2, Parent: 1, AuthorName: "Bob", DateGmt: "2026-06-05T12:27:41", Content: struct {
			Rendered string `json:"rendered"`
		}{Rendered: "Reply"}},
		{ID: 3, Parent: 2, AuthorName: "Charlie", DateGmt: "2026-06-05T12:28:41", Content: struct {
			Rendered string `json:"rendered"`
		}{Rendered: "Nested"}},
	}

	posts := toModel(comments)
	assert.Len(t, posts, 3)
	assert.Equal(t, 0, posts[0].Depth)
	assert.Equal(t, 1, posts[1].Depth)
	assert.Equal(t, 2, posts[2].Depth)
}

func TestToModelMixedOrder(t *testing.T) {
	comments := []WordPressComment{
		{ID: 3, Parent: 1, AuthorName: "Charlie", DateGmt: "2026-06-05T12:28:41", Content: struct {
			Rendered string `json:"rendered"`
		}{Rendered: "Nested reply"}},
		{ID: 1, Parent: 0, AuthorName: "Alice", DateGmt: "2026-06-05T12:26:41", Content: struct {
			Rendered string `json:"rendered"`
		}{Rendered: "Top level"}},
		{ID: 2, Parent: 1, AuthorName: "Bob", DateGmt: "2026-06-05T12:27:41", Content: struct {
			Rendered string `json:"rendered"`
		}{Rendered: "Direct reply"}},
	}

	posts := toModel(comments)
	assert.Len(t, posts, 3)
	assert.Equal(t, 0, posts[1].Depth, "parent should be depth 0")
	assert.Equal(t, 1, posts[0].Depth, "nested reply should be depth 1")
	assert.Equal(t, 1, posts[2].Depth, "direct reply should be depth 1")
}

func TestToModelEmpty(t *testing.T) {
	posts := toModel(nil)
	assert.Empty(t, posts)

	posts = toModel([]WordPressComment{})
	assert.Empty(t, posts)
}
