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
	var action string
	action, a.activeThread, a.activePost = ui.HandleViewerInput(a.threads, a.activeThread, a.activePost)
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
			go func() {
				a.ContinueStub(ui)
			}()
		}
	case "quit":
		a.run = false
	}
}

func (a *App) ContinueStub(ui ui.UI) {
	activeThread := &a.threads[a.activeThread]
	activePostIndex := a.activePost

	activePost := activeThread.Posts[activePostIndex]
	if activePost.Stub.Key == "" {
		a.Error("No more replies can be fetched on this thread")
		return
	}

	a.SiteInput.ContinueFrom = &model.ContinueFrom{
		Key:   activePost.Stub.Key,
		Depth: activePost.Depth,
	}
	posts, err := a.Site.Fetch(a.SiteInput)
	if err != nil {
		a.Error(err.Error())
		return
	}

	activeThread.Posts = activeThread.Posts.
		RemoveAt(activePostIndex). // remove stub
		AppendAt(posts, activePostIndex)

	if len(posts) < activePost.Stub.Count {
		activeThread.Posts = append(activeThread.Posts, model.Post{
			ID:    uuid.NewString(),
			Depth: activePost.Depth,
			Stub: &model.Stub{
				Count: activePost.Stub.Count - len(posts),
				Key:   "", // TODO
			},
		})
	}

	ui.Refresh()
}
