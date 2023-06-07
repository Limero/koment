package reddit

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/limero/koment/lib/helper"
	"github.com/limero/koment/lib/model"
)

type Reddit struct{}

func NewReddit() Reddit {
	return Reddit{}
}

func (s Reddit) GetInput(url *url.URL, _ ...string) (*model.SiteInput, error) {
	if !strings.Contains(url.Path, "/comments/") {
		return nil, fmt.Errorf("invalid path %q", url.Path)
	}
	return &model.SiteInput{
		SiteName: model.SiteReddit,
		Category: strings.Split(strings.Split(url.Path, "/r/")[1], "/")[0],
		ID:       strings.Split(strings.Split(url.Path, "/comments/")[1], "/")[0],
	}, nil
}

func (s Reddit) Fetch(fi model.SiteInput) (model.Posts, error) {
	if fi.Demo {
		return s.getFromExampleFile()
	}
	return s.getFromApi(fi.Category, fi.ID)
}

func (s Reddit) getFromApi(subReddit string, threadID string) (model.Posts, error) {
	var resp Listings
	if err := helper.GetPageToJSON(fmt.Sprintf(
		"https://reddit.com/r/%s/comments/%s.json",
		subReddit,
		threadID,
	), &resp); err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, fmt.Errorf("no posts found, probably rate limited")
	}

	return resp.toModel()
}

func (s Reddit) getFromExampleFile() (model.Posts, error) {
	data, err := os.ReadFile("lib/demo/reddit-example-api.json")
	if err != nil {
		return nil, err
	}

	var resp Listings
	if err = json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return resp.toModel()
}
