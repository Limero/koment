package app

import (
	"errors"
	"fmt"

	"github.com/limero/koment/app/info"
	"github.com/limero/koment/app/ui"
	"github.com/limero/koment/lib"
	"github.com/limero/koment/lib/model"
)

type App struct {
	Site      model.Site
	SiteInput model.SiteInput

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
		a.accumulatedPosts = append(a.accumulatedPosts, posts...)
		a.threads = append(a.threads, model.Thread{Posts: posts})
		ui.Refresh()
	}

	a.loadingMsg = fmt.Sprintf("Loading comments from %s...", a.SiteInput.SiteName)
	ui.DrawLoading(a.loadingMsg)

	go func() {
		posts, err := a.Site.Fetch(a.SiteInput)
		a.loadingMsg = ""

		if err != nil {
			a.Fatal(err.Error())
			ui.Refresh()
			return
		}

		if len(a.accumulatedPosts) > 0 {
			// Already populated incrementally via OnPost (hackernews)
		} else {
			a.threads = model.PostsToThreads(posts)
			if len(posts) == 0 {
				a.Terminate("No comments available")
			}
		}

		ui.Refresh()
	}()

	for a.run {
		if len(a.threads) > 0 {
			ui.DrawViewer(
				a.threads,
				a.activeThread,
				a.activePost,
			)
		} else if a.loadingMsg != "" {
			ui.DrawLoading(a.loadingMsg)
		}

		if a.infoMsg != "" {
			ui.DrawInfo(a.infoLevel, a.infoMsg)
			if a.infoLevel == info.InfoLevelFatal || a.infoLevel == info.InfoLevelTerminate {
				ui.PauseUntilInput()
				if a.infoLevel == info.InfoLevelFatal {
					return errors.New(a.infoMsg)
				}
				return nil
			}
			a.infoMsg = ""
		}

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
