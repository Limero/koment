package disqus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToModel(t *testing.T) {
	listPosts := ListPostsThreaded{
		Response: []Post{
			{
				Dislikes:  321,
				Likes:     123,
				ID:        "1",
				CreatedAt: "2023-04-21T13:37:00",
				Author: Author{
					Name: "author",
				},
				RawMessage: "body",
				Depth:      0,
			},
			{
				ID:         "2",
				CreatedAt:  "2023-04-21T13:37:00",
				RawMessage: "replybody",
				Depth:      1,
			},
		},
	}

	posts, err := listPosts.toModel()
	require.NoError(t, err)
	assert.Len(t, posts, 2)

	t.Run("Check post", func(t *testing.T) {
		expected := listPosts.Response[0]
		actual := posts[0]
		assert.Equal(t, expected.ID, actual.ID)
		assert.Equal(t, expected.RawMessage, actual.Message)
		assert.Equal(t, expected.Depth, actual.Depth)
		assert.Equal(t, expected.Author.Name, actual.Author.Name)
		assert.Equal(t, expected.CreatedAt, actual.CreatedAt.Format("2006-01-02T15:04:05"))
		assert.Equal(t, expected.Likes, *actual.Upvotes)
		assert.Equal(t, expected.Dislikes, *actual.Downvotes)
	})

	t.Run("Check reply", func(t *testing.T) {
		expected := listPosts.Response[1]
		actual := posts[1]
		assert.Equal(t, expected.ID, actual.ID)
		assert.Equal(t, expected.RawMessage, actual.Message)
		assert.Equal(t, expected.Depth, actual.Depth)
	})
}
