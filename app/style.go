package app

import "github.com/gdamore/tcell/v2"

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

	// Command/Search
	Command tcell.Style

	// Info message
	LoadingMessage tcell.Style
	InfoMessage    tcell.Style
	ErrorMessage   tcell.Style
}

func DefaultStyle() Style {
	d := tcell.StyleDefault
	return Style{
		MessageLength: 60,
		FullIndent:    4,
		SemiIndent:    2,

		AuthorName:    d.Foreground(tcell.ColorBlue),
		UpVotesIcon:   d.Foreground(tcell.ColorGray),
		UpVotesNum:    d.Foreground(tcell.ColorGray),
		DownVotesIcon: d.Foreground(tcell.ColorGray),
		DownVotesNum:  d.Foreground(tcell.ColorGray),
		Time:          d.Foreground(tcell.ColorGray),
		Separator:     d.Foreground(tcell.ColorGray),

		AuthorStartChar: '▎',
		UpVotesChar:     tcell.RuneUArrow,
		DownVotesChar:   tcell.RuneDArrow,
		SeparatorChar:   tcell.RuneBullet,

		ActiveMessage:     d.Foreground(tcell.ColorOrange),
		RegularMessage:    d,
		StubMessage:       d.Foreground(tcell.ColorGray),
		ActiveStubMessage: d.Foreground(tcell.ColorOrange),

		ActiveMessageChar: tcell.RuneVLine,
		StubStartChar:     '▎',

		Command: d,

		LoadingMessage: d,
		InfoMessage:    d,
		ErrorMessage:   d.Background(tcell.ColorRed),
	}
}
