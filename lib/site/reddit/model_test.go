package reddit

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToModel(t *testing.T) {
	replies := Listing{
		Data: ListingData{
			Children: []Children{
				{
					Data: ChildrenData{
						ID:         "2",
						Body:       "replybody",
						Depth:      1,
						CreatedUTC: json.Number("1680813139.0"),
					},
				},
			},
		},
	}
	repliesJson, err := json.Marshal(replies)
	require.NoError(t, err)

	listings := Listings{
		{}, // first index is metadata on post
		{
			Data: ListingData{
				Children: []Children{
					{
						Data: ChildrenData{
							ID:         "1",
							Body:       "body",
							Depth:      0,
							Author:     "author",
							CreatedUTC: json.Number("1680813139.0"),
							Ups:        123,
							Downs:      321,
							RepliesRaw: repliesJson,
						},
					},
				},
			},
		},
	}

	posts, err := listings.toModel()
	require.NoError(t, err)
	assert.Len(t, posts, 2)

	t.Run("Check post", func(t *testing.T) {
		expected := listings[1].Data.Children[0].Data
		actual := posts[0]
		assert.Equal(t, expected.ID, actual.ID)
		assert.Equal(t, expected.Body, actual.Message)
		assert.Equal(t, expected.Depth, actual.Depth)
		assert.Equal(t, expected.Author, actual.Author.Name)
		assert.Equal(t, expected.CreatedUTC.String(), fmt.Sprintf("%d.0", actual.CreatedAt.Unix()))
		assert.Equal(t, expected.Ups, *actual.Upvotes)
		assert.Equal(t, expected.Downs, *actual.Downvotes)
	})

	t.Run("Check reply", func(t *testing.T) {
		expected := replies.Data.Children[0].Data
		actual := posts[1]
		assert.Equal(t, expected.ID, actual.ID)
		assert.Equal(t, expected.Body, actual.Message)
		assert.Equal(t, expected.Depth, actual.Depth)
	})
}
