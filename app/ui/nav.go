package ui

import (
	"github.com/gdamore/tcell/v3"
	"github.com/limero/koment/app/util"
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
			case "h":
				return "hide-post", t, p
			case "n":
				return "search-next", t, p
			case "N":
				return "search-prev", t, p
		case "/":
				return "search", t, p
			}
		case tcell.KeyCtrlL:
			ui.screen.Sync()
		case tcell.KeyCtrlC:
			return "quit", t, p
		case tcell.KeyCtrlD:
			t, p = ui.navDownHalfPage(threads, t, p)
		case tcell.KeyCtrlU:
			t, p = ui.navUpHalfPage(threads, t, p)
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

func (ui *ui) HandleSearchInput() (string, string) {
	ev := <-ui.screen.EventQ()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		ui.Refresh()
		switch ev.Key() {
		case tcell.KeyRune:
			return "search-add", ev.Str()
		case tcell.KeyBackspace:
			return "search-rm", ""
		case tcell.KeyEnter:
			return "search-exec", ""
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

func (ui *ui) countPostLines(thread model.Thread, p int) int {
	post := thread.Posts[p]
	if post.Stub != nil {
		return 1
	}
	if hasHiddenParent(thread, p) {
		return 0
	}
	if post.Hidden {
		return 1
	}
	return 1 + len(util.TextToLines(post.Message, ui.style.MessageLength))
}

func (ui *ui) navDownHalfPage(threads model.Threads, t, p int) (int, int) {
	_, screenH := ui.screen.Size()
	targetLines := max(1, screenH/2)
	lines := 0

	for lines < targetLines {
		if p+1 < len(threads[t].Posts) {
			p++
		} else if t+1 < len(threads) {
			t++
			p = 0
			lines++ // blank line between threads
		} else {
			break
		}
		lines += ui.countPostLines(threads[t], p)
	}

	return t, p
}

func (ui *ui) navUpHalfPage(threads model.Threads, t, p int) (int, int) {
	_, screenH := ui.screen.Size()
	targetLines := max(1, screenH/2)
	lines := 0

	for lines < targetLines {
		if p > 0 {
			p--
		} else if t > 0 {
			t--
			p = len(threads[t].Posts) - 1
			lines++ // blank line between threads
		} else {
			break
		}
		lines += ui.countPostLines(threads[t], p)
	}

	return t, p
}

func hasHiddenParent(thread model.Thread, p int) bool {
	depth := thread.Posts[p].Depth
	for i := 0; i < p; i++ {
		if thread.Posts[i].Hidden && thread.Posts[i].Depth < depth {
			return true
		}
	}
	return false
}

func navUpPost(threads model.Threads, t, p int) (int, int) {
	p--
	if p < 0 && t > 0 {
		t--
		p = len(threads[t].Posts) - 1
	}
	for p >= 0 && hasHiddenParent(threads[t], p) {
		p--
		if p < 0 && t > 0 {
			t--
			p = len(threads[t].Posts) - 1
		} else if p < 0 {
			break
		}
	}
	return t, max(0, p)
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
	for hasHiddenParent(threads[t], p) {
		p++
		if p >= len(threads[t].Posts) {
			if t < len(threads)-1 {
				t++
				p = 0
			} else {
				break
			}
		}
	}

	return t, p
}

func navUpThread(t int) (int, int) {
	return max(0, t-1), 0
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
