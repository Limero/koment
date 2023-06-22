package ycombinator

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToModel(t *testing.T) {
	listPosts := Posts{
		{
			By:    "author",
			ID:    1,
			Text:  "body",
			Score: 123,
			Time:  1680813139,
		},
	}

	depth := 0
	posts, err := listPosts.toModel()
	require.NoError(t, err)
	assert.Len(t, posts, 1)

	t.Run("Check post", func(t *testing.T) {
		expected := listPosts[0]
		actual := posts[0]
		assert.Equal(t, strconv.Itoa(expected.ID), actual.ID)
		assert.Equal(t, expected.Text, actual.Message)
		assert.Equal(t, depth, actual.Depth)
		assert.Equal(t, expected.By, actual.Author.Name)
		assert.Equal(t, expected.Time, actual.CreatedAt.Unix())
		assert.Equal(t, expected.Score, *actual.Upvotes)
		assert.Nil(t, actual.Downvotes)
	})
}
