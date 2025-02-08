package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/limero/koment/app/info"
	"github.com/limero/koment/lib/model"
)

type UI interface {
	// Draw
	DrawLoading(msg string)
	DrawViewer(threads model.Threads, activePostID string)
	DrawCommandPrompt(command string)
	DrawInfo(infoLevel info.InfoLevel, msg string)
	Refresh()

	// Nav
	HandleViewerInput(threads model.Threads, t, p int) (string, int, int)
	HandleCommandInput() (string, rune)
	PauseUntilInput()
}

type ui struct {
	screen tcell.Screen
	view   *views.ViewPort
	style  Style

	shouldCenter bool
}

func New(style Style) (*ui, error) {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)

	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, fmt.Errorf("error creating screen: %w", err)
	}

	if err = screen.Init(); err != nil {
		return nil, fmt.Errorf("error initializing screen: %w", err)
	}

	return &ui{
		screen: screen,
		style:  style,
		view:   views.NewViewPort(screen, 0, 0, -1, -1),
	}, nil
}

func (ui ui) Fini() {
	ui.screen.Fini()
}
