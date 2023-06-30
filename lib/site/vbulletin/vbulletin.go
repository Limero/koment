package vbulletin

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/limero/koment/lib/model"
)

type Vbulletin struct{}

func NewVbulletin() Vbulletin {
	return Vbulletin{}
}

func (s Vbulletin) GetInput(url *url.URL, _ ...string) (*model.SiteInput, error) {
	if !strings.Contains(url.Path, "/forum/") {
		return nil, fmt.Errorf("invalid path %q", url.Path)
	}

	return &model.SiteInput{
		SiteName: model.SiteVbulletin,
		FullUrl:  url,
	}, nil
}

func (s Vbulletin) Fetch(fi model.SiteInput) (model.Posts, error) {
	return s.getFromHttp(fi.FullUrl)
}

func (s Vbulletin) getFromHttp(url *url.URL) (model.Posts, error) {
	res, err := http.Get(url.String())
	if err != nil {
		return model.Posts{}, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return model.Posts{}, err
	}

	fromPosts := doc.Find(".b-post")
	posts := make(model.Posts, fromPosts.Length())
	fromPosts.Each(func(i int, s *goquery.Selection) {
		createdAtInt, _ := strconv.ParseInt(s.AttrOr("data-node-publishdate", ""), 10, 64)
		createdAt := time.Unix(createdAtInt, 0)

		upvotes, _ := strconv.Atoi(s.Find(".votecount").Text())

		posts[i] = model.Post{
			ID:    s.AttrOr("data-node-id", ""),
			Depth: 0,
			Author: model.Author{
				Name: s.Find(".author a").Text(),
			},
			Message: s.Find(".js-post__content-text").Text(),

			Upvotes:   &upvotes,
			CreatedAt: &createdAt,
		}
	})

	return posts, nil
}
