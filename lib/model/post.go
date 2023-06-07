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
