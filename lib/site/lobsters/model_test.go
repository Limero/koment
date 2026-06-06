package lobsters

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToModel(t *testing.T) {
	story := Story{
		ShortID:       "abc123",
		Title:         "Test Story",
		SubmitterUser: "op_user",
		Comments: []Comment{
			{
				ShortID:        "c1",
				CreatedAt:      "2023-01-01T00:00:00.000Z",
				Score:          42,
				ParentComment:  nil,
				Depth:          0,
				CommentPlain:   "A top-level comment",
				Comment:        "<p>A top-level comment</p>",
				CommentingUser: "user1",
				IsDeleted:      false,
				IsModerated:    false,
			},
			{
				ShortID:        "c2",
				CreatedAt:      "2023-01-01T01:00:00.000Z",
				Score:          7,
				ParentComment:  strPtr("c1"),
				Depth:          1,
				CommentPlain:   "A reply",
				Comment:        "<p>A reply</p>",
				CommentingUser: "op_user",
				IsDeleted:      false,
				IsModerated:    false,
			},
			{
				ShortID:        "c3",
				CreatedAt:      "2023-01-01T02:00:00.000Z",
				Score:          -3,
				ParentComment:  nil,
				Depth:          0,
				CommentPlain:   "",
				Comment:        "<p>No plain text</p>",
				CommentingUser: "user2",
				IsDeleted:      true,
				IsModerated:    false,
			},
		},
	}

	posts, err := story.toModel()
	require.NoError(t, err)
	require.Len(t, posts, 3)

	t.Run("Top-level comment", func(t *testing.T) {
		p := posts[0]
		assert.Equal(t, "c1", p.ID)
		assert.Equal(t, 0, p.Depth)
		assert.Equal(t, "user1", p.Author.Name)
		assert.Equal(t, "A top-level comment", p.Message)
		assert.False(t, p.Hidden)
		assert.False(t, p.IsOP)
		assert.Equal(t, 42, *p.Upvotes)
		assert.Equal(t, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), *p.CreatedAt)
	})

	t.Run("Reply from OP", func(t *testing.T) {
		p := posts[1]
		assert.Equal(t, "c2", p.ID)
		assert.Equal(t, 1, p.Depth)
		assert.Equal(t, "op_user", p.Author.Name)
		assert.Equal(t, "A reply", p.Message)
		assert.False(t, p.Hidden)
		assert.True(t, p.IsOP)
		assert.Equal(t, 7, *p.Upvotes)
	})

	t.Run("Deleted comment falls back to HTML", func(t *testing.T) {
		p := posts[2]
		assert.Equal(t, "c3", p.ID)
		assert.Equal(t, 0, p.Depth)
		assert.Equal(t, "user2", p.Author.Name)
		assert.Equal(t, "<p>No plain text</p>", p.Message)
		assert.True(t, p.Hidden)
		assert.Equal(t, -3, *p.Upvotes)
	})
}

// TODO: When bumping to Go 1.26, use new() instead of this
func strPtr(s string) *string {
	return &s
}
