package main

import (
	"fmt"
	"log"
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
		log.Fatal(err)
	}

	app := app.NewApp()
	app.SiteInput = *siteInput

	if err := app.InitScreen(); err != nil {
		log.Fatal(err)
	}
	if err := app.RunApp(); err != nil {
		log.Fatal(err)
	}
}
