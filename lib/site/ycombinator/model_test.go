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
			Kids:  []int{2, 3},
			Text:  "body",
			Score: 123,
			Time:  1680813139,
		},
	}

	depth := 0
	posts, err := listPosts.toModel(depth)
	require.NoError(t, err)
	assert.Len(t, posts, 3)

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

	t.Run("Check replies stubs", func(t *testing.T) {
		expected := listPosts[0]
		actual := posts[1]
		require.NotNil(t, actual.Stub)
		assert.Len(t, actual.ID, 36) // generated uuid
		assert.Equal(t, depth+1, actual.Depth)
		assert.Equal(t, 1, actual.Stub.Count)
		assert.Equal(t, strconv.Itoa(expected.Kids[0]), actual.Stub.Key)

		actual = posts[2]
		require.NotNil(t, actual.Stub)
		assert.Len(t, actual.ID, 36) // generated uuid
		assert.Equal(t, depth+1, actual.Depth)
		assert.Equal(t, 1, actual.Stub.Count)
		assert.Equal(t, strconv.Itoa(expected.Kids[1]), actual.Stub.Key)
	})
}
