package app

import (
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
	Demo      bool

	screen       tcell.Screen
	threads      model.Threads
	search       Search
	activeThread int
	activePost   int
	mode         string
	command      string
	infoMsg      string
	infoLevel    string
	run          bool
}

func NewApp() App {
	return App{
		Style:     DefaultStyle(),
		mode:      "viewer",
		infoLevel: "info",
		run:       true,
	}
}

func (a *App) RunApp() {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	var err error
	a.screen, err = tcell.NewScreen()
	if err != nil {
		fmt.Printf("Error creating screen: %s\n", err)
		return
	}
	defer a.screen.Fini()

	if err = a.screen.Init(); err != nil {
		fmt.Printf("Error initializing screen: %s\n", err)
		return
	}

	a.screen.Clear()
	view := views.NewViewPort(a.screen, 0, 0, -1, -1)

	drawLoading(a.Style, view, fmt.Sprintf("Loading comments from %s...", a.SiteInput.SiteName))
	a.screen.Sync()

	a.SiteInput.Demo = a.Demo
	a.Site = lib.NewSite(a.SiteInput.SiteName)

	posts, err := a.Site.Fetch(a.SiteInput)
	if err != nil {
		a.Fatal(err.Error())
	} else if len(posts) == 0 {
		a.Fatal("No comments available")
	}

	a.threads = model.PostsToThreads(posts)

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
			if a.infoLevel == "fatal" {
				a.screen.Sync()
				PauseUntilInput(a.screen)
				// TODO: No need to panic, but just println doesn't seem to output anything
				panic(a.infoMsg)
			}
			a.infoMsg = ""
		}

		switch a.mode {
		case "command":
			drawCommandPrompt(a.Style, view, a.command)
			a.CommandMode()
		case "viewer":
			a.ViewerMode()
		}
	}
}
