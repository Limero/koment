package app

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2/views"
	"github.com/limero/koment/lib/model"
)

func drawLoading(style Style, view *views.ViewPort, msg string) {
	for x, c := range []rune(msg) {
		view.SetContent(x, 0, c, nil, style.LoadingMessage)
	}
}

func drawAuthorLine(style Style, view *views.ViewPort, post model.Post, x int, y int) int {
	view.SetContent(x, y, style.AuthorStartChar, nil, style.AuthorStart)
	x++
	for _, c := range []rune(post.Author.Name) {
		view.SetContent(x, y, c, nil, style.AuthorName)
		x++
	}
	if post.Upvotes != nil || post.Downvotes != nil {
		x++
		view.SetContent(x, y, style.SeparatorChar, nil, style.Separator)
		x++
	}
	if post.Upvotes != nil {
		x++
		view.SetContent(x, y, style.UpVotesChar, nil, style.UpVotesIcon)
		x++
		for _, c := range []rune(strconv.Itoa(*post.Upvotes)) {
			view.SetContent(x, y, c, nil, style.UpVotesNum)
			x++
		}
	}
	if post.Downvotes != nil {
		x++
		view.SetContent(x, y, style.DownVotesChar, nil, style.DownVotesIcon)
		x++
		for _, c := range []rune(strconv.Itoa(*post.Downvotes)) {
			view.SetContent(x, y, c, nil, style.DownVotesNum)
			x++
		}
	}
	if post.CreatedAt != nil {
		x++
		view.SetContent(x, y, style.SeparatorChar, nil, style.Separator)
		x += 2
		for _, c := range []rune(post.CreatedAt.Format(time.DateTime)) {
			view.SetContent(x, y, c, nil, style.Time)
			x++
		}
	}

	return x
}

func drawStub(style Style, view *views.ViewPort, post model.Post, activePostID string, x int, y int) {
	st := style.StubMessage
	if post.ID == activePostID {
		x++
		st = style.ActiveStubMessage
	}
	view.SetContent(x, y, style.StubStartChar, nil, st)
	x++
	for _, c := range []rune(fmt.Sprintf("%d more replies", post.Stub.Count)) {
		view.SetContent(x, y, c, nil, st)
		x++
	}
}

func drawViewer(style Style, view *views.ViewPort, threads model.Threads, activePostID string) (int, int) {
	x, y := 0, 0
	activeMsgLength := 0
	activeMsgY := 0

	for _, thread := range threads {
		for _, post := range thread.Posts {
			x = post.Depth * style.FullIndent

			if post.Stub != nil {
				drawStub(style, view, post, activePostID, x, y)
				if post.ID == activePostID {
					activeMsgY = y
					activeMsgLength = 1
				}
				y++
				continue
			}

			// Author line
			x = drawAuthorLine(style, view, post, x, y)
			y++

			// Main message
			lines := TextToLines(post.Message, style.MessageLength)
			if post.ID == activePostID {
				activeMsgLength = len(lines)
				activeMsgY = y
			}
			for _, line := range lines {
				x = (post.Depth * style.FullIndent) + style.SemiIndent
				if post.ID == activePostID {
					view.SetContent(x, y, style.ActiveMessageChar, nil, style.ActiveMessage)
					x++
				}

				for _, c := range []rune(line) {
					view.SetContent(x, y, c, nil, style.RegularMessage)
					x++
				}
				y++
			}
		}
		y++
	}

	return activeMsgLength, activeMsgY
}

func drawCommandPrompt(style Style, view *views.ViewPort, command string) {
	_, _, width, height := view.GetVisible()

	x := 0
	for _, c := range []rune(":" + command) {
		view.SetContent(x, height, c, nil, style.Command)
		x++
	}
	for ; x <= width; x++ {
		view.SetContent(x, height, ' ', nil, style.Command)
	}
}

func drawInfo(style Style, view *views.ViewPort, infoLevel InfoLevel, msg string) {
	_, _, width, height := view.GetVisible()

	messageStyle := style.InfoMessage
	switch infoLevel {
	case InfoLevelError, InfoLevelFatal:
		messageStyle = style.ErrorMessage
	}

	x := 0
	for _, c := range []rune(msg) {
		view.SetContent(x, height, c, nil, messageStyle)
		x++
	}
	for ; x <= width; x++ {
		view.SetContent(x, height, ' ', nil, messageStyle)
	}
}
