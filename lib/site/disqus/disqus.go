package disqus

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/limero/koment/lib/internal/util"
	"github.com/limero/koment/lib/model"
)

type Disqus struct{}

func NewDisqus() Disqus {
	return Disqus{}
}

func (s Disqus) GetInput(url *url.URL, v ...string) (*model.SiteInput, error) {
	number, err := util.GetNumberFromPath(url.Path)
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
		return nil, errors.Join(err, errors.New("failed to retrieve disqus api key"))
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

func (s Disqus) getFromApi(apiKey string, threadID string) (model.Posts, error) {
	limit := 100
	order := "popular"
	cursor := "1%3A0%3A0"

	var resp ListPostsThreaded
	if err := util.GetPageToJSON(fmt.Sprintf(
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
		return "", fmt.Errorf("either name (%s) or number (%s) is empty", name, number)
	}
	url := "https://disqus.com/embed/comments/?f=" + name + "&t_i=" + number + "#version=93621f724643ecd0f307feb8123718cb"

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	var data EmbedPage

	jsonData := doc.Find("script#disqus-threadData").First().Text()
	if err = json.Unmarshal([]byte(jsonData), &data); err != nil {
		return "", err
	}
	return data.Response.Thread.ID, nil
}
