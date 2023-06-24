package disqus

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/limero/koment/lib/helper"
	"github.com/limero/koment/lib/model"
)

type Disqus struct{}

func NewDisqus() Disqus {
	return Disqus{}
}

func (s Disqus) GetInput(url *url.URL, v ...string) (*model.SiteInput, error) {
	number, err := helper.GetNumberFromPath(url.Path)
	if err != nil {
		return nil, err
	}
	if len(v) == 0 {
		return nil, errors.New("Disqus requires additional variables to decide input")
	}
	threadID, err := s.getThreadIDFromEmbedPage(v[0], number)
	if err != nil {
		return nil, err
	}

	apiKey, err := s.getApiKey()
	if err != nil {
		return nil, err
	}

	return &model.SiteInput{
		SiteName: model.SiteDisqus,
		ID:       threadID,
		ApiKey:   apiKey,
	}, nil
}

func (s Disqus) Fetch(fi model.SiteInput) (model.Posts, error) {
	return s.getFromApi(fi.ApiKey, fi.ID)
}

func (s Disqus) getApiKey() (string, error) {
	apiKeyFile := helper.CachePath("disqus-api-key.txt")
	apiKey, err := helper.ReadFileIfExists(apiKeyFile)
	if err != nil {
		return "", err
	}
	if apiKey != "" {
		return apiKey, nil
	}

	url := "https://c.disquscdn.com/next/embed/common.bundle.6719fe9dbe70a5a047052a905ea1cbc5.js"
	body, err := helper.GetPageBodyString(url)
	if err != nil {
		return "", err
	}
	apiKey, err = helper.GetLastBetween(body, "embedAPI:\"", "\"")
	if err != nil {
		return "", err
	}

	return apiKey, helper.WriteFile(apiKey, apiKeyFile)
}

func (s Disqus) getFromApi(apiKey string, threadID string) (model.Posts, error) {
	limit := 100
	order := "popular"
	cursor := "1%3A0%3A0"

	var resp ListPostsThreaded
	if err := helper.GetPageToJSON(fmt.Sprintf(
		"https://disqus.com/api/3.0/threads/listPostsThreaded?limit=%d&thread=%s&order=%s&cursor=%s&api_key=%s",
		limit,
		threadID,
		order,
		cursor,
		apiKey,
	), &resp); err != nil {
		return nil, err
	}

	return resp.toModel()
}

func (s Disqus) getThreadIDFromEmbedPage(name string, number string) (string, error) {
	if name == "" || number == "" {
		return "", fmt.Errorf("Either name (%s) or number (%s) is empty", name, number)
	}
	url := "https://disqus.com/embed/comments/?f=" + name + "&t_i=" + number + "#version=93621f724643ecd0f307feb8123718cb"
	body, err := helper.GetPageBodyString(url)
	if err != nil {
		return "", err
	}
	return helper.GetLastBetween(body, "\"id\":\"", "\"")
}
