package main

import (
	"fmt"
	"os"

	"github.com/limero/koment/app"
	"github.com/limero/koment/app/ui"
	"github.com/limero/koment/lib"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Usage: koment <url>")
		return
	}

	siteInput, err := lib.FindComments(args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	app := app.NewApp()
	app.SiteInput = *siteInput

	style := ui.DefaultStyle()
	ui, err := ui.New(style)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer ui.Fini()

	if err := app.RunApp(ui); err != nil {
		ui.Fini()
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
