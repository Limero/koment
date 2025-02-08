package app

import (
	"testing"

	"github.com/limero/koment/app/test"
	"github.com/limero/koment/lib/model"
)

func TestRunApp(t *testing.T) {
	a := NewApp()
	a.SiteInput = model.SiteInput{
		SiteName: model.SiteDemo,
	}
	ui := test.MockUI{}

	go func() {
		_ = a.RunApp(&ui)
	}()
}
