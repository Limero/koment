package app

import "github.com/limero/koment/lib/model"

type Search struct {
	Term    string
	Results model.Posts
	Index   int
}

func (a *App) SearchStart(term string) {
	a.mu.Lock()
	a.search = Search{
		Term:    term,
		Results: a.threads.FindPostsContaining(term),
		Index:   0,
	}
	if len(a.search.Results) > 0 {
		a.activeThread, a.activePost = a.threads.FindPost(a.search.Results[a.search.Index].ID)
	}
	a.mu.Unlock()
	a.Info("Found %d result(s)", len(a.search.Results))
}

func (a *App) SearchNext() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.search.Results) == 0 {
		return
	}
	if a.search.Index < len(a.search.Results)-1 {
		a.search.Index++
	}
	a.activeThread, a.activePost = a.threads.FindPost(a.search.Results[a.search.Index].ID)
}

func (a *App) SearchPrev() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.search.Results) == 0 {
		return
	}
	if a.search.Index > 0 {
		a.search.Index--
	}
	a.activeThread, a.activePost = a.threads.FindPost(a.search.Results[a.search.Index].ID)
}
