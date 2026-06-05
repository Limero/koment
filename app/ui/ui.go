package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v3"
	"github.com/limero/koment/app/info"
	"github.com/limero/koment/lib/model"
)

type UI interface {
	// Draw
	DrawLoading(msg string)
	DrawViewer(threads model.Threads, activeThread, activePost int)
	DrawCommandPrompt(command string)
	DrawInfo(infoLevel info.InfoLevel, msg string)
	Refresh()

	// Nav
	HandleViewerInput(threads model.Threads, t, p int) (string, int, int)
	HandleCommandInput() (string, string)
	PauseUntilInput()
}

type ui struct {
	screen tcell.Screen
	style  Style

	shouldCenter bool
	scrollY      int
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
	}, nil
}

func (ui ui) Fini() {
	ui.screen.Fini()
}
