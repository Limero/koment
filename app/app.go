package app

import (
	"errors"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/limero/koment/lib"
	"github.com/limero/koment/lib/model"
)

type App struct {
	Site      model.Site
	SiteInput model.SiteInput
	Style     Style

	screen       tcell.Screen
	threads      model.Threads
	search       Search
	activeThread int
	activePost   int
	mode         Mode
	command      string
	infoMsg      string
	infoLevel    InfoLevel
	run          bool
	loading      bool
}

func NewApp() App {
	return App{
		Style:     DefaultStyle(),
		mode:      ModeViewer,
		infoLevel: InfoLevelInfo,
		run:       true,
	}
}

func (a *App) InitScreen() error {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	var err error
	a.screen, err = tcell.NewScreen()
	if err != nil {
		return fmt.Errorf("Error creating screen: %w", err)
	}

	if err = a.screen.Init(); err != nil {
		return fmt.Errorf("Error initializing screen: %w", err)
	}

	return nil
}

func (a *App) RunApp() error {
	defer a.screen.Fini()
	view := views.NewViewPort(a.screen, 0, 0, -1, -1)

	a.Site = lib.NewSite(a.SiteInput.SiteName)

	a.loading = true
	go func() {
		posts, err := a.Site.Fetch(a.SiteInput)
		if err != nil {
			a.Fatal(err.Error())
		} else if len(posts) == 0 {
			a.Terminate("No comments available")
		}

		a.threads = model.PostsToThreads(posts)
		a.loading = false
		a.Refresh()
	}()

	shouldCenter := false
	for a.run {
		a.screen.Clear()

		if len(a.threads) > 0 {
			activePostLines, activePostY := drawViewer(
				a.Style,
				view,
				a.threads,
				a.threads[a.activeThread].Posts[a.activePost].ID,
			)

			// TODO: Find a better way to center the message instead of doing this every time
			if shouldCenter {
				shouldCenter = false
				view.Center(0, activePostY+(activePostLines/2))
				continue
			}
			shouldCenter = true
		}

		if a.infoMsg != "" {
			drawInfo(a.Style, view, a.infoLevel, a.infoMsg)
			if a.infoLevel == InfoLevelFatal || a.infoLevel == InfoLevelTerminate {
				a.screen.Show()
				PauseUntilInput(a.screen)
				if a.infoLevel == InfoLevelFatal {
					return errors.New(a.infoMsg)
				}
				return nil
			}
			a.infoMsg = ""
		}

		if a.loading {
			drawLoading(a.Style, view, fmt.Sprintf("Loading comments from %s...", a.SiteInput.SiteName))
			a.screen.Show()
		}

		switch a.mode {
		case ModeCommand:
			drawCommandPrompt(a.Style, view, a.command)
			a.CommandMode()
		case ModeViewer:
			a.ViewerMode()
		}
	}
	return nil
}

func (a *App) Refresh() {
	if a.screen == nil {
		return
	}
	// Send a nil event to the waiting event listener to redraw everything
	_ = a.screen.PostEvent(tcell.NewEventInterrupt(nil))
}
