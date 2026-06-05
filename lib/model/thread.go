package model

import "strings"

type Thread struct {
	Posts Posts
}

type Threads []Thread

func PostsToThreads(posts Posts) Threads {
	numThreads := 0
	for _, p := range posts {
		if p.Depth == 0 {
			numThreads++
		}
	}
	threads := make(Threads, 0, numThreads)

	var thread *Thread
	for i, p := range posts {
		if p.Depth == 0 {
			if thread != nil {
				threads = append(threads, *thread)
			}
			thread = &Thread{}
		}
		if thread == nil {
			thread = &Thread{}
		}
		thread.Posts = append(thread.Posts, posts[i])
	}
	if thread != nil {
		threads = append(threads, *thread)
	}
	return threads
}

func (threads Threads) FindPostsContaining(s string) Posts {
	posts := make(Posts, 0)
	for _, thread := range threads {
		for i := range thread.Posts {
			if strings.Contains(thread.Posts[i].Message, s) {
				posts = append(posts, thread.Posts[i])
			}
		}
	}
	return posts
}

func (threads Threads) FindPost(postID string) (int, int) {
	for threadIndex, thread := range threads {
		for postIndex, post := range thread.Posts {
			if post.ID == postID {
				return threadIndex, postIndex
			}
		}
	}
	return 0, 0
}

func (threads Threads) TotalPosts() int {
	total := 0
	for _, t := range threads {
		total += len(t.Posts)
	}
	return total
}

func (threads Threads) CurrentPostIndex(activeThread, activePost int) int {
	index := 0
	for t := 0; t < activeThread; t++ {
		index += len(threads[t].Posts)
	}
	return index + activePost + 1 // 1-indexed
}
