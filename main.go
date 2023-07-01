package main

import (
	"fmt"
	"os"

	"github.com/limero/koment/app"
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
		fmt.Println("Error:", err)
		return
	}

	app := app.NewApp()
	app.SiteInput = *siteInput
	app.RunApp()
}
