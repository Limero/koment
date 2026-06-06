package ui

import (
	"fmt"
	"hash/fnv"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
	"github.com/limero/koment/app/info"
	"github.com/limero/koment/app/util"
	"github.com/limero/koment/lib/model"
)

var authorPalette = []tcell.Color{
	color.Orange,
	color.Violet,
	color.Brown,
	color.Indigo,
	color.Coral,
	color.Crimson,
	color.DarkCyan,
	color.DarkOrange,
	color.DarkOrchid,
	color.DarkSeaGreen,
	color.DeepPink,
	color.DodgerBlue,
	color.ForestGreen,
	color.Goldenrod,
	color.HotPink,
	color.IndianRed,
	color.MediumPurple,
	color.RebeccaPurple,
	color.RoyalBlue,
	color.SeaGreen,
	color.SlateBlue,
	color.SteelBlue,
	color.Tomato,
	color.Turquoise,
	color.YellowGreen,
}

func (ui *ui) colorForAuthor(name string) tcell.Color {
	if c, ok := ui.authorColors[name]; ok {
		return c
	}
	h := fnv.New32a()
	h.Write([]byte(name))
	c := authorPalette[h.Sum32()%uint32(len(authorPalette))]
	ui.authorColors[name] = c
	return c
}

func (ui *ui) DrawLoading(msg string) {
	ui.screen.Clear()
	for x, c := range []rune(msg) {
		ui.setContent(x, 0, c, nil, ui.style.LoadingMessage)
	}
	ui.screen.Show()
}

func (ui *ui) DrawViewer(threads model.Threads, activeThread, activePost int) {
	x, y := 0, 0
	activeMsgLength := 0
	activeMsgY := 0
	activePostID := threads[activeThread].Posts[activePost].ID

	ui.screen.Clear()
	for _, thread := range threads {
		for pi := range thread.Posts {
			post := &thread.Posts[pi]
			x = post.Depth * ui.style.FullIndent

			if post.Stub != nil {
				ui.drawStub(*post, activePostID, x, y)
				if post.ID == activePostID {
					activeMsgY = y
					activeMsgLength = 1
				}
				y++
				continue
			}

			if hasHiddenParent(thread, pi) {
				if post.ID == activePostID {
					activeMsgY = y
					activeMsgLength = 1
				}
				continue
			}

			// Author line
			ui.drawAuthorLine(*post, activePostID, x, y)
			y++

			// Main message
			if post.Hidden {
				if post.ID == activePostID {
					activeMsgY = y
					activeMsgLength = 1
				}
				y++
				continue
			}
			lines := util.TextToLines(post.Message, ui.style.MessageLength)
			if post.ID == activePostID {
				activeMsgLength = len(lines)
				activeMsgY = y
			}
			for _, line := range lines {
				x = (post.Depth * ui.style.FullIndent) + ui.style.SemiIndent
				if post.ID == activePostID {
					ui.setContent(x, y, ui.style.ActiveMessageChar, nil, ui.style.ActiveMessage)
					x++
				}

				for _, c := range line {
					ui.setContent(x, y, c, nil, ui.style.RegularMessage)
					x++
				}
				y++
			}
		}
		y++
	}

	screenW, screenH := ui.screen.Size()
	if ui.shouldCenter {
		ui.shouldCenter = false
		ui.scrollY = (activeMsgY + activeMsgLength/2) - screenH/2
		if ui.scrollY < 0 {
			ui.scrollY = 0
		}
		return
	}
	ui.shouldCenter = true

	totalPosts := threads.TotalPosts()
	if totalPosts > 0 {
		currentPos := threads.CurrentPostIndex(activeThread, activePost)
		progressStr := fmt.Sprintf(" %d/%d ", currentPos, totalPosts)
		px := screenW - len(progressStr)
		for i, c := range progressStr {
			ui.screen.SetContent(px+i, screenH-1, c, nil, ui.style.Progress)
		}
	}

	ui.screen.Show()
}

func (ui *ui) setContent(x, y int, ch rune, comb []rune, style tcell.Style) {
	ui.screen.SetContent(x, y-ui.scrollY, ch, comb, style)
}

func (ui *ui) drawAuthorLine(post model.Post, activePostID string, x int, y int) {
	authorStartStyle := ui.style.AuthorStart
	if post.Hidden {
		authorStartStyle = tcell.StyleDefault.Foreground(color.Gray)
	}
	if post.ID == activePostID {
		authorStartStyle = ui.style.ActiveMessage
	}
	ui.setContent(x, y, ui.style.AuthorStartChar, nil, authorStartStyle)
	x++
	authorColor := ui.colorForAuthor(post.Author.Name)
	for _, c := range post.Author.Name {
		ui.setContent(x, y, c, nil, tcell.StyleDefault.Foreground(authorColor))
		x++
	}
	if post.IsOP {
		x++
		ui.setContent(x, y, '[', nil, ui.style.OPBadge)
		x++
		ui.setContent(x, y, 'O', nil, ui.style.OPBadge)
		x++
		ui.setContent(x, y, 'P', nil, ui.style.OPBadge)
		x++
		ui.setContent(x, y, ']', nil, ui.style.OPBadge)
	}
	if post.Upvotes != nil || post.Downvotes != nil {
		x++
		ui.setContent(x, y, ui.style.SeparatorChar, nil, ui.style.Separator)
		x++
	}
	if post.Upvotes != nil {
		x++
		ui.setContent(x, y, ui.style.UpVotesChar, nil, ui.style.UpVotesIcon)
		x++
		for _, c := range strconv.Itoa(*post.Upvotes) {
			ui.setContent(x, y, c, nil, ui.style.UpVotesNum)
			x++
		}
	}
	if post.Downvotes != nil {
		x++
		ui.setContent(x, y, ui.style.DownVotesChar, nil, ui.style.DownVotesIcon)
		x++
		for _, c := range strconv.Itoa(*post.Downvotes) {
			ui.setContent(x, y, c, nil, ui.style.DownVotesNum)
			x++
		}
	}
	if post.CreatedAt != nil {
		x++
		ui.setContent(x, y, ui.style.SeparatorChar, nil, ui.style.Separator)
		x += 2
		for _, c := range post.CreatedAt.Format(time.DateTime) {
			ui.setContent(x, y, c, nil, ui.style.Time)
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
	ui.setContent(x, y, ui.style.StubStartChar, nil, st)
	x++
	for _, c := range fmt.Sprintf("%d more replies", post.Stub.Count) {
		ui.setContent(x, y, c, nil, st)
		x++
	}
}

func (ui *ui) DrawSearchPrompt(search string) {
	width, height := ui.screen.Size()

	x := 0
	for _, c := range "/" + search {
		ui.screen.SetContent(x, height-1, c, nil, ui.style.Search)
		x++
	}
	for ; x <= width; x++ {
		ui.screen.SetContent(x, height-1, ' ', nil, ui.style.Search)
	}

	ui.screen.Show()
}

func (ui *ui) DrawInfo(infoLevel info.InfoLevel, msg string) {
	width, height := ui.screen.Size()

	messageStyle := ui.style.InfoMessage
	switch infoLevel {
	case info.InfoLevelError, info.InfoLevelFatal:
		messageStyle = ui.style.ErrorMessage
	}

	x := 0
	for _, c := range msg {
		ui.screen.SetContent(x, height-1, c, nil, messageStyle)
		x++
	}
	for ; x <= width; x++ {
		ui.screen.SetContent(x, height-1, ' ', nil, messageStyle)
	}

	ui.screen.Show()
}

func (ui *ui) Refresh() {
	if ui.stopped.Load() {
		return
	}

	ui.finiMu.Lock()
	defer ui.finiMu.Unlock()
	if ui.stopped.Load() {
		return
	}

	select {
	case ui.screen.EventQ() <- tcell.NewEventInterrupt(nil):
	default:
	}
}
