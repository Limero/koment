package youtube

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/limero/koment/lib/helper"
	"github.com/limero/koment/lib/model"
)

type Youtube struct {
	invidiousInstance string
}

func NewYoutube() Youtube {
	return Youtube{
		invidiousInstance: "https://invidious.snopyta.org",
	}
}

func (s Youtube) GetInput(url *url.URL, _ ...string) (*model.SiteInput, error) {
	queryValues := url.Query()
	return &model.SiteInput{
		SiteName: model.SiteYoutube,
		ID:       queryValues.Get("v"),
	}, nil
}

func (s Youtube) Fetch(fi model.SiteInput) (model.Posts, error) {
	if fi.Demo {
		return s.getFromExampleFile()
	}
	return s.getFromApi(fi.ID, fi.ContinueFrom)
}

func (s Youtube) getFromApi(videoID string, continueFrom string) (model.Posts, error) {
	var resp CommentsResponse
	if err := helper.GetPageToJSON(fmt.Sprintf(
		"%s/api/v1/comments/%s/?continuation=%s",
		s.invidiousInstance,
		videoID,
		continueFrom,
	), &resp); err != nil {
		return nil, err
	}

	depth := 0
	if continueFrom != "" {
		depth = 1
	}

	return resp.toModel(depth)
}

func (s Youtube) getFromExampleFile() (model.Posts, error) {
	data, err := os.ReadFile("lib/demo/youtube-example-api.json")
	if err != nil {
		return nil, err
	}

	var resp CommentsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return resp.toModel(0)
}
