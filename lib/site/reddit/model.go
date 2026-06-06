package reddit

import (
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/limero/koment/lib/internal/util"
	"github.com/limero/koment/lib/model"
)

func parseComments(doc *goquery.Document) (model.Posts, error) {
	opName := doc.Find("div.thing.link").First().AttrOr("data-author", "")

	var posts model.Posts
	doc.Find("div.thing.comment").Each(func(_ int, s *goquery.Selection) {
		if s.AttrOr("data-fullname", "") == "" {
			return
		}
		if s.AttrOr("data-author", "") == "" {
			return
		}
		posts = append(posts, parseComment(s, opName))
	})

	return posts, nil
}

func parseComment(s *goquery.Selection, opName string) model.Post {
	id := s.AttrOr("data-fullname", "")
	author := s.AttrOr("data-author", "")

	depth := s.ParentsFiltered(".thing.comment").Length()

	upvotesStr := s.Find(".score.unvoted").First().AttrOr("title", "")
	upvotes, _ := strconv.Atoi(upvotesStr)

	timeStr := s.Find("time.live-timestamp").First().AttrOr("datetime", "")
	var createdAt time.Time
	if timeStr != "" {
		createdAt, _ = time.Parse(time.RFC3339, timeStr)
	}

	body := ""
	mdSel := s.Find(".usertext-body .md").First()
	if mdSel.Length() > 0 {
		bodyHTML, _ := mdSel.Html()
		body = util.CleanHTML(bodyHTML)
	}
	if body == "" {
		body = strings.TrimSpace(s.Find(".usertext-body").First().Text())
	}

	return model.Post{
		ID:    id,
		Depth: depth,
		Author: model.Author{
			Name: author,
		},
		Message:   body,
		IsOP:      author == opName,
		Upvotes:   &upvotes,
		CreatedAt: &createdAt,
	}
}
