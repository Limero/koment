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

	threads      model.Threads
	search       Search
	activeThread int
	activePost   int
	mode         Mode
	command      string
	infoMsg      string
	infoLevel    info.InfoLevel
	run          bool
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

	ui.DrawLoading(fmt.Sprintf("Loading comments from %s...", a.SiteInput.SiteName))

	go func() {
		posts, err := a.Site.Fetch(a.SiteInput)
		if err != nil {
			a.Fatal(err.Error())
		} else if len(posts) == 0 {
			a.Terminate("No comments available")
		}

		a.threads = model.PostsToThreads(posts)
		ui.Refresh()
	}()

	for a.run {
		if len(a.threads) > 0 {
			ui.DrawViewer(
				a.threads,
				a.threads[a.activeThread].Posts[a.activePost].ID,
			)
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
