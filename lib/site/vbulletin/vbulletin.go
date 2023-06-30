package vbulletin

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/limero/koment/lib/helper"
	"github.com/limero/koment/lib/model"
)

type Vbulletin struct{}

func NewVbulletin() Vbulletin {
	return Vbulletin{}
}

func (s Vbulletin) GetInput(url *url.URL, _ ...string) (*model.SiteInput, error) {
	if !strings.Contains(url.Path, "/forum/") {
		res, err := http.Get(url.String())
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			return nil, err
		}

		var foundUrl string
		doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
			href := s.AttrOr("href", "")
			if strings.Contains(href, "/node/") {
				foundUrl = href
				return
			}
		})
		if foundUrl != "" {
			u, err := url.Parse(foundUrl)
			if err != nil {
				return nil, err
			}
			return &model.SiteInput{
				SiteName: model.SiteVbulletin,
				FullUrl:  u,
			}, nil
		}

		return nil, fmt.Errorf("could not find any comments for %q", url)
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

		// remove any quoted messages
		s.Find(".bbcode_container").Remove()
		s.Find(".b-bbcode").Remove()

		posts[i] = model.Post{
			ID:    s.AttrOr("data-node-id", ""),
			Depth: 0,
			Author: model.Author{
				Name: s.Find(".author a").Text(),
			},
			Message: helper.CleanText(s.Find(".js-post__content-text").Text()),

			Upvotes:   &upvotes,
			CreatedAt: &createdAt,
		}
	})

	return posts, nil
}
