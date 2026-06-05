package ui

import (
	"github.com/gdamore/tcell/v3"
	"github.com/limero/koment/lib/model"
)

func (ui *ui) HandleViewerInput(threads model.Threads, t, p int) (string, int, int) {
	ev := <-ui.screen.EventQ()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		ui.Refresh()
		switch ev.Key() {
		case tcell.KeyRune:
			switch ev.Str() {
			case "q", "Q":
				return "quit", t, p
			case "j":
				t, p = navDownPost(threads, t, p)
			case "J":
				t, p = navDownThread(threads, t)
			case "k":
				t, p = navUpPost(threads, t, p)
			case "K":
				t, p = navUpThread(t)
			case "g":
				t, p = navTop()
			case "G":
				t, p = navBottom(threads)
			case "n":
				return "search-next", t, p
			case "N":
				return "search-prev", t, p
			case ":":
				return "command", t, p
			case "/":
				return "search", t, p
			}
		case tcell.KeyCtrlL:
			ui.screen.Sync()
		case tcell.KeyCtrlC:
			return "quit", t, p
		case tcell.KeyUp:
			if ev.Modifiers() == tcell.ModShift {
				t, p = navUpThread(t)
			} else {
				t, p = navUpPost(threads, t, p)
			}
		case tcell.KeyDown:
			if ev.Modifiers() == tcell.ModShift {
				t, p = navDownThread(threads, t)
			} else {
				t, p = navDownPost(threads, t, p)
			}
		case tcell.KeyEnter:
			return "enter", t, p
		}
	case *tcell.EventResize:
		ui.screen.Sync()
	case *tcell.EventMouse:
		ui.Refresh()
		switch ev.Buttons() {
		case tcell.WheelUp:
			t, p = navUpPost(threads, t, p)
		case tcell.WheelDown:
			t, p = navDownPost(threads, t, p)
		}
	}

	return "", t, p
}

func (ui *ui) HandleCommandInput() (string, string) {
	ev := <-ui.screen.EventQ()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		ui.Refresh()
		switch ev.Key() {
		case tcell.KeyRune:
			return "command-add", ev.Str()
		case tcell.KeyBackspace:
			return "command-rm", ""
		case tcell.KeyEnter:
			return "command-exec", ""
		case tcell.KeyESC:
			return "exit", ""
		case tcell.KeyCtrlC:
			return "quit", ""
		}
	case *tcell.EventResize:
		ui.screen.Sync()
	}

	return "", ""
}

func (ui *ui) PauseUntilInput() {
	for {
		ev := <-ui.screen.EventQ()
		switch ev.(type) {
		case *tcell.EventKey:
			return
		}
	}
}

func navUpPost(threads model.Threads, t, p int) (int, int) {
	p--
	if p < 0 {
		if t > 0 {
			t--
			p = len(threads[t].Posts) - 1
		} else {
			p = 0
		}
	}

	return t, p
}

func navDownPost(threads model.Threads, t, p int) (int, int) {
	p++
	if p >= len(threads[t].Posts) {
		if t < len(threads)-1 {
			t++
			p = 0
		} else {
			p--
		}
	}

	return t, p
}

func navUpThread(t int) (int, int) {
	if t > 0 {
		t--
	}
	return t, 0
}

func navDownThread(threads model.Threads, t int) (int, int) {
	p := 0
	if t < len(threads)-1 {
		t++
	} else {
		p = len(threads[t].Posts) - 1
	}
	return t, p
}

func navTop() (int, int) {
	return 0, 0
}

func navBottom(threads model.Threads) (int, int) {
	t := len(threads) - 1
	return t, len(threads[t].Posts) - 1
}
