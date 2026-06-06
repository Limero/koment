package app

import "github.com/limero/koment/app/ui"

func (a *App) SetSearchMode() {
	a.searchInput = ""
	a.mode = ModeSearch
}

func (a *App) SearchMode(ui ui.UI) {
	action, char := ui.HandleSearchInput()
	switch action {
	case "search-add":
		a.searchInput += char
	case "search-rm":
		if len(a.searchInput) > 0 {
			a.searchInput = a.searchInput[:len(a.searchInput)-1]
		}
	case "search-exec":
		a.SetViewerMode()
		a.SearchStart(a.searchInput)
	case "exit":
		a.searchInput = ""
		a.SetViewerMode()
	case "quit":
		a.run = false
	}
}
