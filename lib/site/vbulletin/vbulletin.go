package vbulletin

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/limero/koment/lib/internal/util"
	"github.com/limero/koment/lib/model"
)

type Vbulletin struct{}

func NewVbulletin() Vbulletin {
	return Vbulletin{}
}

func (s Vbulletin) GetInput(u *url.URL, _ ...string) (*model.SiteInput, error) {
	if !strings.Contains(u.Path, "/forum/") {
		res, err := util.HTTPGet(u.String())
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			return nil, err
		}

		var foundURL string
		doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
			href := s.AttrOr("href", "")
			if strings.Contains(href, "/node/") {
				foundURL = href
				return
			}
		})
		if foundURL != "" {
			u, err := url.Parse(foundURL)
			if err != nil {
				return nil, err
			}
			return &model.SiteInput{
				SiteName: model.SiteVbulletin,
				FullURL:  u,
			}, nil
		}

		return nil, fmt.Errorf("could not find any comments for %q", u)
	}

	return &model.SiteInput{
		SiteName: model.SiteVbulletin,
		FullURL:  u,
	}, nil
}

func (s Vbulletin) Fetch(fi model.SiteInput) (model.Posts, error) {
	return s.getFromHTTP(fi.FullURL)
}

func (s Vbulletin) getFromHTTP(url *url.URL) (model.Posts, error) {
	res, err := util.HTTPGet(url.String())
	if err != nil {
		return model.Posts{}, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return model.Posts{}, err
	}

	type rawPost struct {
		post    model.Post
		replyTo string
	}

	order := make([]string, 0)
	allPosts := make(map[string]rawPost)
	firstAuthor := ""
	doc.Find(".b-post").Each(func(_ int, s *goquery.Selection) {
		createdAtInt, _ := strconv.ParseInt(s.AttrOr("data-node-publishdate", ""), 10, 64)
		createdAt := time.Unix(createdAtInt, 0)

		upvotes, _ := strconv.Atoi(s.Find(".votecount").Text())

		replyTo := s.Find("a[title='View Post']").AttrOr("href", "")
		replyToURL, _ := url.Parse(replyTo)
		replyTo = replyToURL.Query().Get("p")

		// remove any quoted messages
		s.Find(".bbcode_container").Remove()
		s.Find(".b-bbcode").Remove()

		author := s.Find(".author a").Text()
		if firstAuthor == "" {
			firstAuthor = author
		}

		id := s.AttrOr("data-node-id", "")
		order = append(order, id)
		allPosts[id] = rawPost{
			post: model.Post{
				ID:    id,
				Depth: 0,
				Author: model.Author{
					Name: author,
				},
				Message: util.CleanHTML(s.Find(".js-post__content-text").Text()),
				IsOP:    author == firstAuthor,

				Upvotes:   &upvotes,
				CreatedAt: &createdAt,
			},
			replyTo: replyTo,
		}
	})

	children := make(map[string][]string)
	roots := make([]string, 0)
	for _, id := range order {
		rp := allPosts[id]
		if rp.replyTo != "" {
			if _, ok := allPosts[rp.replyTo]; ok {
				children[rp.replyTo] = append(children[rp.replyTo], id)
			} else {
				roots = append(roots, id)
			}
		} else {
			roots = append(roots, id)
		}
	}

	posts := make(model.Posts, 0, len(allPosts))
	var dfs func(id string, depth int)
	dfs = func(id string, depth int) {
		rp := allPosts[id]
		rp.post.Depth = depth
		posts = append(posts, rp.post)
		for _, cid := range children[id] {
			dfs(cid, depth+1)
		}
	}
	for _, rootID := range roots {
		dfs(rootID, 0)
	}

	return posts, nil
}
