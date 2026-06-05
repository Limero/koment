package ycombinator

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/limero/koment/lib/internal/util"
	"github.com/limero/koment/lib/model"
)

type Ycombinator struct{}

func NewYcombinator() Ycombinator {
	return Ycombinator{}
}

func (s Ycombinator) GetInput(url *url.URL, _ ...string) (*model.SiteInput, error) {
	queryValues := url.Query()
	return &model.SiteInput{
		SiteName: model.SiteYcombinator,
		ID:       queryValues.Get("id"),
	}, nil
}

func (s Ycombinator) Fetch(fi model.SiteInput) (model.Posts, error) {
	return s.getFromApi(fi)
}

func (s Ycombinator) getFromApi(fi model.SiteInput) (model.Posts, error) {
	id := fi.ID
	depth := 0
	if fi.ContinueFrom != nil {
		id = fi.ContinueFrom.Key
		depth = fi.ContinueFrom.Depth
	}

	url := "https://hacker-news.firebaseio.com/v0/item/%s.json"
	var resp Post
	if err := util.GetPageToJSON(fmt.Sprintf(url, id), &resp); err != nil {
		return nil, err
	}

	var posts model.Posts
	switch resp.Type {
	case "story":
		for _, kid := range resp.Kids {
			var newResp Post
			if err := util.GetPageToJSON(fmt.Sprintf(url, strconv.Itoa(kid)), &newResp); err != nil {
				return nil, err
			}

			batch := newResp.toModelBatch(depth)
			posts = append(posts, batch...)

			if fi.OnPost != nil {
				fi.OnPost(batch)
			}
		}
	case "comment":
		posts = append(posts, resp.toModelBatch(depth)...)
	}

	return posts, nil
}
