package app

import (
	"strings"

	"github.com/limero/koment/app/ui"
)

func (a *App) SetCommandMode(cmd string) {
	a.command = cmd
	a.mode = ModeCommand
}

func (a *App) CommandMode(ui ui.UI) {
	action, char := ui.HandleCommandInput()
	switch action {
	case "command-add":
		a.command += string(char)
	case "command-rm":
		if len(a.command) > 0 {
			a.command = a.command[:len(a.command)-1]
		}
	case "command-exec":
		a.ExecCommand(a.command)
	case "exit":
		a.SetViewerMode()
	case "quit":
		a.run = false
	}
}

func (a *App) ExecCommand(cmd string) {
	a.SetViewerMode()
	if strings.HasPrefix(cmd, "search ") {
		term, _ := strings.CutPrefix(cmd, "search ")
		a.SearchStart(term)
		return
	}

	if strings.ToLower(a.command) == "q" {
		a.run = false
		return
	}

	a.Error("Invalid command: %q", a.command)
}
