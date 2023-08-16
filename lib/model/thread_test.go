package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostsToThreads(t *testing.T) {
	author := Author{
		Name: "Name",
	}

	posts := Posts{
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
			Posts: Posts{
				posts[0],
				posts[1],
			},
		},
		{
			Posts: Posts{
				posts[2],
			},
		},
	}

	assert.Equal(t, expected, PostsToThreads(posts))
}

func TestFindPostsContaining(t *testing.T) {
	threads := Threads{
		{
			Posts: Posts{
				{Message: "lorem ipsum"},
			},
		},
		{
			Posts: Posts{
				{Message: "hello"},
				{Message: "lorem ipsum"},
			},
		},
	}

	for _, tt := range []struct {
		term          string
		expectedPosts Posts
	}{
		{
			term: "ore",
			expectedPosts: Posts{
				threads[0].Posts[0],
				threads[1].Posts[1],
			},
		},
		{
			term: "e",
			expectedPosts: Posts{
				threads[0].Posts[0],
				threads[1].Posts[0],
				threads[1].Posts[1],
			},
		},
		{
			term: "hello",
			expectedPosts: Posts{
				threads[1].Posts[0],
			},
		},
		{
			term:          "HELLO",
			expectedPosts: Posts{},
		},
	} {
		posts := threads.FindPostsContaining(tt.term)
		assert.Equal(t, tt.expectedPosts, posts)
	}
}

func TestFindPost(t *testing.T) {
	threads := Threads{
		{
			Posts: Posts{
				{ID: "1"},
				{ID: "2"},
			},
		},
		{
			Posts: Posts{
				{ID: "3"},
				{ID: "4"},
			},
		},
	}

	for _, tt := range []struct {
		postID         string
		expectedThread int
		expectedPost   int
	}{
		{
			postID:         "1",
			expectedThread: 0,
			expectedPost:   0,
		},
		{
			postID:         "2",
			expectedThread: 0,
			expectedPost:   1,
		},
		{
			postID:         "3",
			expectedThread: 1,
			expectedPost:   0,
		},
		{
			postID:         "4",
			expectedThread: 1,
			expectedPost:   1,
		},
		{
			postID:         "not-found",
			expectedThread: 0,
			expectedPost:   0,
		},
	} {
		threadIndex, postIndex := threads.FindPost(tt.postID)
		assert.Equal(t, tt.expectedThread, threadIndex)
		assert.Equal(t, tt.expectedPost, postIndex)
	}
}
