package model

import (
	"slices"
	"sort"
	"time"
)

type Author struct {
	Name string
}

type Post struct {
	ID      string
	Depth   int
	Author  Author
	Message string

	Upvotes   *int
	Downvotes *int
	CreatedAt *time.Time

	Stub *Stub
}

type Posts []Post

func (posts Posts) RemoveAt(index int) Posts {
	if index >= 0 && index < len(posts) {
		return slices.Delete(posts, index, index+1)
	}
	return posts
}

func (posts Posts) AppendAt(newPosts Posts, index int) Posts {
	if index < 0 || index >= len(posts) {
		return append(posts, newPosts...)
	}
	return slices.Insert(posts, index, newPosts...)
}

func (posts Posts) SortByDepth() {
	sort.Slice(posts, func(i, j int) bool {
		if posts[i].Depth == 0 || posts[j].Depth == 0 {
			return false
		}
		return posts[i].Depth < posts[j].Depth
	})
}
