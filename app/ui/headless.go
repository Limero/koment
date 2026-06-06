package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/limero/koment/app/util"
	"github.com/limero/koment/lib/model"
)

type HeadlessUI struct {
	Style Style
}

func NewHeadless(style Style) HeadlessUI {
	return HeadlessUI{Style: style}
}

func (h HeadlessUI) Render(threads model.Threads) {
	for ti, thread := range threads {
		if ti > 0 {
			fmt.Println()
		}
		for pi := range thread.Posts {
			post := &thread.Posts[pi]
			indent := strings.Repeat(" ", post.Depth*h.Style.FullIndent)

			if post.Stub != nil {
				fmt.Printf("%s\u2590 %d more replies\n", indent, post.Stub.Count)
				continue
			}

			if hasHiddenParent(thread, pi) {
				continue
			}

			opBadge := ""
			if post.IsOP {
				opBadge = " [OP]"
			}
			fmt.Printf("%s\u2590 %s%s", indent, post.Author.Name, opBadge)
			if post.Upvotes != nil {
				fmt.Printf(" \u2191%d", *post.Upvotes)
			}
			if post.Downvotes != nil {
				fmt.Printf(" \u2193%d", *post.Downvotes)
			}
			if post.CreatedAt != nil {
				fmt.Printf(" \u2022 %s", post.CreatedAt.Format(time.DateTime))
			}
			fmt.Println()

			if post.Hidden {
				fmt.Printf("%s  [hidden]\n", indent)
				continue
			}

			msgIndent := indent + strings.Repeat(" ", h.Style.SemiIndent)
			lines := util.TextToLines(post.Message, h.Style.MessageLength)
			for _, line := range lines {
				fmt.Printf("%s%s\n", msgIndent, line)
			}

			if pi < len(thread.Posts)-1 {
				fmt.Println()
			}
		}
	}
}

func (h HeadlessUI) RenderPosts(posts model.Posts) {
	threads := model.PostsToThreads(posts)
	h.Render(threads)
}
