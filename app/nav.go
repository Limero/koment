package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/limero/koment/lib/model"
)

func HandleViewerInput(screen tcell.Screen, threads model.Threads, t, p int) (string, int, int) {
	ev := screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyRune:
			switch ev.Rune() {
			case 'q', 'Q':
				return "quit", t, p
			case 'j':
				t, p = NavDownPost(threads, t, p)
			case 'J':
				t, p = NavDownThread(threads, t)
			case 'k':
				t, p = NavUpPost(threads, t, p)
			case 'K':
				t, p = NavUpThread(t)
			case 'g':
				t, p = NavTop()
			case 'G':
				t, p = NavBottom(threads)
			case 'n':
				return "search-next", t, p
			case 'N':
				return "search-prev", t, p
			case ':':
				return "command", t, p
			case '/':
				return "search", t, p
			}
		case tcell.KeyCtrlL:
			screen.Sync()
		case tcell.KeyCtrlC:
			return "quit", t, p
		case tcell.KeyUp:
			if ev.Modifiers() == tcell.ModShift {
				t, p = NavUpThread(t)
			} else {
				t, p = NavUpPost(threads, t, p)
			}
		case tcell.KeyDown:
			if ev.Modifiers() == tcell.ModShift {
				t, p = NavDownThread(threads, t)
			} else {
				t, p = NavDownPost(threads, t, p)
			}
		case tcell.KeyEnter:
			return "enter", t, p
		}
	case *tcell.EventResize:
		screen.Sync()
	case *tcell.EventMouse:
		switch ev.Buttons() {
		case tcell.WheelUp:
			t, p = NavUpPost(threads, t, p)
		case tcell.WheelDown:
			t, p = NavDownPost(threads, t, p)
		}
	}
	return "", t, p
}

func HandleCommandInput(screen tcell.Screen) (string, rune) {
	ev := screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyRune:
			return "command-add", ev.Rune()
		case tcell.KeyBackspace2:
			return "command-rm", 0
		case tcell.KeyEnter:
			return "command-exec", 0
		case tcell.KeyESC:
			return "exit", 0
		case tcell.KeyCtrlC:
			return "quit", 0
		}
	case *tcell.EventResize:
		screen.Sync()
	}
	return "", 0
}

func PauseUntilInput(screen tcell.Screen) {
	for {
		ev := screen.PollEvent()
		switch ev.(type) {
		case *tcell.EventKey:
			return
		}
	}
}

func NavUpPost(threads model.Threads, t, p int) (int, int) {
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

func NavDownPost(threads model.Threads, t, p int) (int, int) {
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

func NavUpThread(t int) (int, int) {
	if t > 0 {
		t--
	}
	return t, 0
}

func NavDownThread(threads model.Threads, t int) (int, int) {
	p := 0
	if t < len(threads)-1 {
		t++
	} else {
		p = len(threads[t].Posts) - 1
	}
	return t, p
}

func NavTop() (int, int) {
	return 0, 0
}

func NavBottom(threads model.Threads) (int, int) {
	t := len(threads) - 1
	return t, len(threads[t].Posts) - 1
}
