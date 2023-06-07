package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeThreads(t *testing.T) {
	author := Author{
		Name: "Name",
	}

	posts := []Post{
		{
			Depth:   0,
			Author:  author,
			Message: "main1",
		},
		{
			Depth:   1,
			Author:  author,
			Message: "nested",
		},
		{
			Depth:   0,
			Author:  author,
			Message: "main2",
		},
	}

	expected := Threads{
		{
			Posts: []Post{
				posts[0],
				posts[1],
			},
		},
		{
			Posts: []Post{
				posts[2],
			},
		},
	}

	assert.Equal(t, expected, PostsToThreads(posts))
}
