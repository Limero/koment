package ui

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/limero/koment/app/helper"
	"github.com/limero/koment/app/info"
	"github.com/limero/koment/lib/model"
)

func (ui *ui) DrawLoading(msg string) {
	ui.screen.Clear()
	for x, c := range []rune(msg) {
		ui.view.SetContent(x, 0, c, nil, ui.style.LoadingMessage)
	}
	ui.screen.Show()
}

func (ui *ui) DrawViewer(threads model.Threads, activePostID string) {
	x, y := 0, 0
	activeMsgLength := 0
	activeMsgY := 0

	ui.screen.Clear()
	for _, thread := range threads {
		for _, post := range thread.Posts {
			x = post.Depth * ui.style.FullIndent

			if post.Stub != nil {
				ui.drawStub(post, activePostID, x, y)
				if post.ID == activePostID {
					activeMsgY = y
					activeMsgLength = 1
				}
				y++
				continue
			}

			// Author line
			ui.drawAuthorLine(post, x, y)
			y++

			// Main message
			lines := helper.TextToLines(post.Message, ui.style.MessageLength)
			if post.ID == activePostID {
				activeMsgLength = len(lines)
				activeMsgY = y
			}
			for _, line := range lines {
				x = (post.Depth * ui.style.FullIndent) + ui.style.SemiIndent
				if post.ID == activePostID {
					ui.view.SetContent(x, y, ui.style.ActiveMessageChar, nil, ui.style.ActiveMessage)
					x++
				}

				for _, c := range line {
					ui.view.SetContent(x, y, c, nil, ui.style.RegularMessage)
					x++
				}
				y++
			}
		}
		y++
	}

	// TODO: Find a better way to center the message instead of doing this every time
	if ui.shouldCenter {
		ui.shouldCenter = false
		ui.view.Center(0, activeMsgY+(activeMsgLength/2))
		return
	}
	ui.shouldCenter = true

	ui.screen.Show()
}

func (ui *ui) drawAuthorLine(post model.Post, x int, y int) {
	ui.view.SetContent(x, y, ui.style.AuthorStartChar, nil, ui.style.AuthorStart)
	x++
	for _, c := range post.Author.Name {
		ui.view.SetContent(x, y, c, nil, ui.style.AuthorName)
		x++
	}
	if post.Upvotes != nil || post.Downvotes != nil {
		x++
		ui.view.SetContent(x, y, ui.style.SeparatorChar, nil, ui.style.Separator)
		x++
	}
	if post.Upvotes != nil {
		x++
		ui.view.SetContent(x, y, ui.style.UpVotesChar, nil, ui.style.UpVotesIcon)
		x++
		for _, c := range strconv.Itoa(*post.Upvotes) {
			ui.view.SetContent(x, y, c, nil, ui.style.UpVotesNum)
			x++
		}
	}
	if post.Downvotes != nil {
		x++
		ui.view.SetContent(x, y, ui.style.DownVotesChar, nil, ui.style.DownVotesIcon)
		x++
		for _, c := range strconv.Itoa(*post.Downvotes) {
			ui.view.SetContent(x, y, c, nil, ui.style.DownVotesNum)
			x++
		}
	}
	if post.CreatedAt != nil {
		x++
		ui.view.SetContent(x, y, ui.style.SeparatorChar, nil, ui.style.Separator)
		x += 2
		for _, c := range post.CreatedAt.Format(time.DateTime) {
			ui.view.SetContent(x, y, c, nil, ui.style.Time)
			x++
		}
	}
}

func (ui *ui) drawStub(post model.Post, activePostID string, x int, y int) {
	st := ui.style.StubMessage
	if post.ID == activePostID {
		x++
		st = ui.style.ActiveStubMessage
	}
	ui.view.SetContent(x, y, ui.style.StubStartChar, nil, st)
	x++
	for _, c := range fmt.Sprintf("%d more replies", post.Stub.Count) {
		ui.view.SetContent(x, y, c, nil, st)
		x++
	}
}

func (ui *ui) DrawCommandPrompt(command string) {
	_, _, width, height := ui.view.GetVisible()

	x := 0
	for _, c := range ":" + command {
		ui.view.SetContent(x, height, c, nil, ui.style.Command)
		x++
	}
	for ; x <= width; x++ {
		ui.view.SetContent(x, height, ' ', nil, ui.style.Command)
	}

	ui.screen.Show()
}

func (ui *ui) DrawInfo(infoLevel info.InfoLevel, msg string) {
	_, _, width, height := ui.view.GetVisible()

	messageStyle := ui.style.InfoMessage
	switch infoLevel {
	case info.InfoLevelError, info.InfoLevelFatal:
		messageStyle = ui.style.ErrorMessage
	}

	x := 0
	for _, c := range msg {
		ui.view.SetContent(x, height, c, nil, messageStyle)
		x++
	}
	for ; x <= width; x++ {
		ui.view.SetContent(x, height, ' ', nil, messageStyle)
	}

	ui.screen.Show()
}

func (ui *ui) Refresh() {
	// Send a nil event to the waiting event listener to redraw everything
	_ = ui.screen.PostEvent(tcell.NewEventInterrupt(nil))
}
