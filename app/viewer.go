package app

import (
	"github.com/google/uuid"
	"github.com/limero/koment/app/ui"
	"github.com/limero/koment/lib/model"
)

func (a *App) SetViewerMode() {
	a.mode = ModeViewer
}

func (a *App) ViewerMode(ui ui.UI) {
	a.mu.Lock()
	if len(a.threads) == 0 {
		a.mu.Unlock()
		return
	}
	threads := a.threads
	at := a.activeThread
	ap := a.activePost
	a.mu.Unlock()

	var action string
	action, at, ap = ui.HandleViewerInput(threads, at, ap)

	a.mu.Lock()
	a.activeThread = at
	a.activePost = ap
	a.mu.Unlock()

	switch action {
	case "search":
		a.SetSearchMode()
	case "search-next":
		a.SearchNext()
	case "search-prev":
		a.SearchPrev()
	case "enter":
		a.mu.Lock()
		if at >= len(a.threads) || ap >= len(a.threads[at].Posts) {
			a.mu.Unlock()
			return
		}
		post := &a.threads[at].Posts[ap]
		if post.Hidden {
			post.Hidden = false
			a.mu.Unlock()
		} else if post.Stub != nil && post.Stub.Key != "" {
			key := post.Stub.Key
			count := post.Stub.Count
			depth := post.Depth
			post.Stub.Key = ""
			a.mu.Unlock()
			go func(k string, d, c int) {
				a.ContinueStub(ui, at, ap, k, d, c)
			}(key, depth, count)
		} else {
			a.mu.Unlock()
		}
	case "hide-post":
		a.mu.Lock()
		if at < len(a.threads) && ap < len(a.threads[at].Posts) {
			post := &a.threads[at].Posts[ap]
			if post.Stub == nil {
				post.Hidden = !post.Hidden
			}
		}
		a.mu.Unlock()
	case "quit":
		a.run = false
	}
}

func (a *App) ContinueStub(ui ui.UI, threadIdx, postIdx int, key string, depth, count int) {
	if key == "" {
		return
	}

	a.mu.Lock()
	if threadIdx >= len(a.threads) || postIdx >= len(a.threads[threadIdx].Posts) {
		a.mu.Unlock()
		return
	}
	fi := a.SiteInput
	fi.ContinueFrom = &model.ContinueFrom{
		Key:   key,
		Depth: depth,
	}
	a.mu.Unlock()

	posts, err := a.Site.Fetch(fi)
	if err != nil {
		a.mu.Lock()
		if threadIdx < len(a.threads) && postIdx < len(a.threads[threadIdx].Posts) {
			p := &a.threads[threadIdx].Posts[postIdx]
			if p.Stub != nil {
				p.Stub.Key = key
			}
		}
		a.mu.Unlock()
		a.Error(err.Error())
		return
	}

	a.mu.Lock()
	if threadIdx >= len(a.threads) || postIdx >= len(a.threads[threadIdx].Posts) {
		a.mu.Unlock()
		return
	}
	activeThread := &a.threads[threadIdx]

	activeThread.Posts = activeThread.Posts.
		RemoveAt(postIdx).
		AppendAt(posts, postIdx)

	if len(posts) < count {
		activeThread.Posts = append(activeThread.Posts, model.Post{
			ID:    uuid.NewString(),
			Depth: depth,
			Stub: &model.Stub{
				Count: count - len(posts),
				Key:   "", // TODO
			},
		})
	}
	a.mu.Unlock()

	ui.Refresh()
	ui.Refresh() // TODO: Shouldn't need to call this twice
}
