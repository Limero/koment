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
	return s.getFromApi(fi.ID, fi.ContinueFrom)
}

func (s Ycombinator) getFromApi(id string, continueFrom *model.ContinueFrom) (model.Posts, error) {
	depth := 0
	if continueFrom != nil {
		id = continueFrom.Key
		depth = continueFrom.Depth
	}

	url := "https://hacker-news.firebaseio.com/v0/item/%s.json"
	var resp Post
	if err := util.GetPageToJSON(fmt.Sprintf(url, id), &resp); err != nil {
		return nil, err
	}

	var posts Posts
	switch resp.Type {
	case "story":
		for _, kid := range resp.Kids {
			var newResp Post
			if err := util.GetPageToJSON(fmt.Sprintf(url, strconv.Itoa(kid)), &newResp); err != nil {
				return nil, err
			}

			posts = append(posts, newResp)
		}
	case "comment":
		posts = append(posts, resp)
	}

	return posts.toModel(depth)
}
