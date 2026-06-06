package app

import (
	"errors"
	"fmt"
	"sync"

	"github.com/limero/koment/app/info"
	"github.com/limero/koment/app/ui"
	"github.com/limero/koment/lib"
	"github.com/limero/koment/lib/model"
)

type App struct {
	Site      model.Site
	SiteInput model.SiteInput

	mu               sync.Mutex
	threads          model.Threads
	search           Search
	activeThread     int
	activePost       int
	mode             Mode
	command          string
	infoMsg          string
	infoLevel        info.InfoLevel
	run              bool
	loadingMsg       string
	accumulatedPosts model.Posts
}

func NewApp() App {
	return App{
		mode:      ModeViewer,
		infoLevel: info.InfoLevelInfo,
		run:       true,
	}
}

func (a *App) RunApp(ui ui.UI) error {
	a.Site = lib.NewSite(a.SiteInput.SiteName)

	a.SiteInput.OnPost = func(posts model.Posts) {
		a.mu.Lock()
		a.accumulatedPosts = append(a.accumulatedPosts, posts...)
		a.threads = append(a.threads, model.Thread{Posts: posts})
		a.mu.Unlock()
		ui.Refresh()
	}

	a.mu.Lock()
	a.loadingMsg = fmt.Sprintf("Loading comments from %s...", a.SiteInput.SiteName)
	a.mu.Unlock()
	ui.DrawLoading(a.loadingMsg)

	go func() {
		posts, err := a.Site.Fetch(a.SiteInput)

		a.mu.Lock()
		a.loadingMsg = ""
		a.mu.Unlock()

		if err != nil {
			a.Fatal(err.Error())
			ui.Refresh()
			return
		}

		a.mu.Lock()
		if len(a.accumulatedPosts) > 0 {
			// Already populated incrementally via OnPost (hackernews)
		} else {
			a.threads = model.PostsToThreads(posts)
			if len(posts) == 0 {
				a.setInfo(info.InfoLevelTerminate, "No comments available")
			}
		}
		a.mu.Unlock()

		ui.Refresh()
	}()

	for a.run {
		a.mu.Lock()
		if len(a.threads) > 0 {
			ui.DrawViewer(a.threads, a.activeThread, a.activePost)
		} else if a.loadingMsg != "" {
			ui.DrawLoading(a.loadingMsg)
		}

		if a.infoMsg != "" {
			im := a.infoMsg
			il := a.infoLevel
			ui.DrawInfo(il, im)
			a.infoMsg = ""
			if il == info.InfoLevelFatal || il == info.InfoLevelTerminate {
				a.mu.Unlock()
				ui.PauseUntilInput()
				if il == info.InfoLevelFatal {
					return errors.New(im)
				}
				return nil
			}
		}
		a.mu.Unlock()

		switch a.mode {
		case ModeCommand:
			ui.DrawCommandPrompt(a.command)
			a.CommandMode(ui)
		case ModeViewer:
			a.ViewerMode(ui)
		}
	}
	return nil
}
