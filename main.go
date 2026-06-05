package main

import (
	"fmt"
	"os"

	"github.com/limero/koment/app"
	"github.com/limero/koment/app/ui"
	"github.com/limero/koment/lib"
)

func main() {
	args := os.Args[1:]

	plain := false
	urlArg := ""

	for _, arg := range args {
		switch arg {
		case "--plain", "-p":
			plain = true
		default:
			urlArg = arg
		}
	}

	if urlArg == "" {
		fmt.Println("Usage: koment [--plain|-p] <url>")
		os.Exit(1)
	}

	siteInput, err := lib.FindComments(urlArg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	a := app.NewApp()
	a.SiteInput = *siteInput

	if plain {
		hu := ui.NewHeadless(ui.DefaultStyle())
		if err := a.RunHeadless(hu); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	style := ui.DefaultStyle()
	tui, err := ui.New(style)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer tui.Fini()

	if err := a.RunApp(tui); err != nil {
		tui.Fini()
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
