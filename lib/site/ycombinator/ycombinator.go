package ycombinator

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/limero/koment/lib/helper"
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
	if fi.Demo {
		return s.getFromExampleFile()
	}
	return s.getFromApi(fi.ID)
}

func (s Ycombinator) getFromApi(id string) (model.Posts, error) {
	url := "https://hacker-news.firebaseio.com/v0/item/%s.json"
	var resp Post
	if err := helper.GetPageToJSON(fmt.Sprintf(url, id), &resp); err != nil {
		return nil, err
	}

	var posts model.Posts

	switch resp.Type {
	case "story":
		for _, kid := range resp.Kids {
			if err := helper.GetPageToJSON(fmt.Sprintf(url, strconv.Itoa(kid)), &resp); err != nil {
				return nil, err
			}

			post, err := resp.toModel()
			if err != nil {
				return nil, err
			}
			posts = append(posts, post)
		}
	case "comment":
		// TODO: Support fetching replies to posts
	}

	return posts, nil
}

func (s Ycombinator) getFromExampleFile() (model.Posts, error) {
	data, err := os.ReadFile("lib/demo/ycombinator-example-api.json")
	if err != nil {
		return nil, err
	}

	var resp Post
	if err = json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	post, err := resp.toModel()
	if err != nil {
		return nil, err
	}

	return model.Posts{post}, nil
}