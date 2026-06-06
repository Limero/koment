package ui

import (
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Style struct {
	MessageLength int
	FullIndent    int
	SemiIndent    int

	// Post author line
	AuthorStart   tcell.Style
	AuthorName    tcell.Style
	UpVotesIcon   tcell.Style
	UpVotesNum    tcell.Style
	DownVotesIcon tcell.Style
	DownVotesNum  tcell.Style
	Time          tcell.Style
	Separator     tcell.Style

	OPBadge tcell.Style

	// Characters
	AuthorStartChar rune
	UpVotesChar     rune
	DownVotesChar   rune
	SeparatorChar   rune

	// Post message
	ActiveMessage     tcell.Style
	RegularMessage    tcell.Style
	ActiveStubMessage tcell.Style
	StubMessage       tcell.Style

	// Characters
	ActiveMessageChar rune
	StubStartChar     rune

	// Search
	Search tcell.Style

	// Info message
	LoadingMessage tcell.Style
	InfoMessage    tcell.Style
	ErrorMessage   tcell.Style

	// Progress
	Progress tcell.Style
}

func DefaultStyle() Style {
	d := tcell.StyleDefault
	return Style{
		MessageLength: 60,
		FullIndent:    4,
		SemiIndent:    2,

		AuthorName:    d.Foreground(color.Blue),
		UpVotesIcon:   d.Foreground(color.Gray),
		UpVotesNum:    d.Foreground(color.Gray),
		DownVotesIcon: d.Foreground(color.Gray),
		DownVotesNum:  d.Foreground(color.Gray),
		Time:          d.Foreground(color.Gray),
		Separator:     d.Foreground(color.Gray),

		OPBadge: d.Foreground(color.Yellow),

		AuthorStartChar: '▎',
		UpVotesChar:     tcell.RuneUArrow,
		DownVotesChar:   tcell.RuneDArrow,
		SeparatorChar:   tcell.RuneBullet,

		ActiveMessage:     d.Foreground(color.Orange),
		RegularMessage:    d,
		StubMessage:       d.Foreground(color.Gray),
		ActiveStubMessage: d.Foreground(color.Orange),

		ActiveMessageChar: tcell.RuneVLine,
		StubStartChar:     '▎',

		Search: d,

		LoadingMessage: d,
		InfoMessage:    d,
		ErrorMessage:   d.Background(color.Red),

		Progress: d.Foreground(color.Gray),
	}
}
