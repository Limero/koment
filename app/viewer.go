package app

import (
	"github.com/google/uuid"
	"github.com/limero/koment/lib/model"
)

func (a *App) SetViewerMode() {
	a.mode = "viewer"
}

func (a *App) ViewerMode() {
	a.screen.Sync()
	var action string
	action, a.activeThread, a.activePost = HandleViewerInput(a.screen, a.threads, a.activeThread, a.activePost)
	switch action {
	case "command":
		a.SetCommandMode("")
	case "search":
		a.SetCommandMode("search ")
	case "search-next":
		a.SearchNext()
	case "search-prev":
		a.SearchPrev()
	case "enter":
		if a.threads[a.activeThread].Posts[a.activePost].Stub != nil {
			a.ContinueStub()
		}
	case "quit":
		a.run = false
	}
}

func (a *App) ContinueStub() {
	activePost := a.threads[a.activeThread].Posts[a.activePost]
	if a.Demo {
		a.Error("Fetching replies does not work in demo")
	} else if activePost.Stub.Key == "" {
		a.Error("No more replies can be fetched on this thread")
	} else {
		a.SiteInput.ContinueFrom = &model.ContinueFrom{
			Key:   activePost.Stub.Key,
			Depth: activePost.Depth,
		}
		posts, err := a.Site.Fetch(a.SiteInput)
		if err != nil {
			a.Error(err.Error())
			return
		}
		a.threads[a.activeThread].Posts = a.threads[a.activeThread].Posts[:len(a.threads[a.activeThread].Posts)-1]
		a.threads[a.activeThread].Posts = append(a.threads[a.activeThread].Posts, posts...)

		if len(posts) < activePost.Stub.Count {
			a.threads[a.activeThread].Posts = append(a.threads[a.activeThread].Posts, model.Post{
				ID:    uuid.NewString(),
				Depth: activePost.Depth,
				Stub: &model.Stub{
					Count: activePost.Stub.Count - len(posts),
					Key:   "", // TODO
				},
			})
		}
	}
}
