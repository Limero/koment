package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveIndex(t *testing.T) {
	posts := Posts{
		{ID: "1"},
		{ID: "2"},
		{ID: "3"},
	}

	t.Run("Remove out of bounds index", func(t *testing.T) {
		newPosts := posts.RemoveAt(100).RemoveAt(-1)
		assert.Equal(t, posts, newPosts)
	})

	t.Run("Remove index", func(t *testing.T) {
		newPosts := posts.RemoveAt(1)
		assert.Len(t, newPosts, 2)
		expectedPosts := Posts{
			{ID: "1"},
			{ID: "3"},
		}
		assert.Equal(t, expectedPosts, newPosts)
	})
}

func TestAppendAt(t *testing.T) {
	posts := Posts{
		{ID: "1"},
		{ID: "2"},
		{ID: "3"},
	}

	t.Run("Append at out of bounds index", func(t *testing.T) {
		newPosts := posts.
			AppendAt(Posts{{ID: "4"}}, 100).
			AppendAt(Posts{{ID: "5"}}, -1)
		assert.Equal(t, Posts{
			{ID: "1"},
			{ID: "2"},
			{ID: "3"},
			{ID: "4"},
			{ID: "5"},
		}, newPosts)
	})

	t.Run("Append at index", func(t *testing.T) {
		newPosts := posts.AppendAt(Posts{
			{ID: "4"},
			{ID: "5"},
		}, 1)
		assert.Len(t, newPosts, 5)
		expectedPosts := Posts{
			{ID: "1"},
			{ID: "4"},
			{ID: "5"},
			{ID: "2"},
			{ID: "3"},
		}
		assert.Equal(t, expectedPosts, newPosts)
	})
}
