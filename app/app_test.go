package app

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/limero/koment/lib/model"
)

func TestRunApp(t *testing.T) {
	a := NewApp()
	a.SiteInput = model.SiteInput{
		SiteName: model.SiteDemo,
	}

	screen := tcell.NewSimulationScreen("")
	a.screen = screen

	go func() {
		_ = a.RunApp()
	}()
}
