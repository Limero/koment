package app

import "strings"

func (a *App) SetCommandMode(cmd string) {
	a.command = cmd
	a.mode = "command"
}

func (a *App) CommandMode() {
	a.screen.Show()
	action, char := HandleCommandInput(a.screen)
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
