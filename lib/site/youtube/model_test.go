package youtube

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToModel(t *testing.T) {
	comments := CommentsResponse{
		Comments: []Comment{
			{
				Author:    "author",
				Content:   "body",
				Published: 1680813139,
				LikeCount: 123,
				CommentID: "1",
				Replies: Replies{
					ReplyCount:   1,
					Continuation: "continuation",
				},
			},
		},
	}

	depth := 0
	posts, err := comments.toModel(depth)
	require.NoError(t, err)
	assert.Len(t, posts, 2)

	t.Run("Check post", func(t *testing.T) {
		expected := comments.Comments[0]
		actual := posts[0]
		assert.Equal(t, expected.CommentID, actual.ID)
		assert.Equal(t, expected.Content, actual.Message)
		assert.Equal(t, depth, actual.Depth)
		assert.Equal(t, expected.Author, actual.Author.Name)
		assert.Equal(t, expected.Published, actual.CreatedAt.Unix())
		assert.Equal(t, expected.LikeCount, *actual.Upvotes)
		assert.Nil(t, actual.Downvotes)
	})

	t.Run("Check replies stub", func(t *testing.T) {
		expected := comments.Comments[0]
		actual := posts[1]
		require.NotNil(t, actual.Stub)
		assert.Len(t, actual.ID, 36) // generated uuid
		assert.Equal(t, depth+1, actual.Depth)
		assert.Equal(t, expected.Replies.ReplyCount, actual.Stub.Count)
		assert.Equal(t, expected.Replies.Continuation, actual.Stub.Key)
	})
}
