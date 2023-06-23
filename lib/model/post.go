package model

import "time"

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
		return append(posts[:index], posts[index+1:]...)
	}
	return posts
}

func (posts Posts) AppendAt(newPosts Posts, index int) Posts {
	if index < 0 || index >= len(posts) {
		return append(posts, newPosts...)
	}

	result := make(Posts, 0, len(posts)+len(newPosts))
	result = append(result, posts[:index]...)
	result = append(result, newPosts...)
	return append(result, posts[index:]...)
}
