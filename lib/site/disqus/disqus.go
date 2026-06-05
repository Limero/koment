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

	return &model.SiteInput{
		SiteName: model.SiteDisqus,
		ID:       number,
		Category: v[0],
	}, nil
}

func (s Disqus) Fetch(fi model.SiteInput) (model.Posts, error) {
	embed, err := s.fetchEmbedPage(fi.Category, fi.ID)
	if err != nil {
		return nil, err
	}
	return embed.toModel()
}

func (s Disqus) fetchEmbedPage(name string, number string) (EmbedPage, error) {
	if name == "" || number == "" {
		return EmbedPage{}, fmt.Errorf("either name (%s) or number (%s) is empty", name, number)
	}
	url := "https://disqus.com/embed/comments/?f=" + name + "&t_i=" + number

	res, err := http.Get(url)
	if err != nil {
		return EmbedPage{}, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return EmbedPage{}, err
	}

	var data EmbedPage

	jsonData := doc.Find("script#disqus-threadData").First().Text()
	if err = json.Unmarshal([]byte(jsonData), &data); err != nil {
		return EmbedPage{}, err
	}
	return data, nil
}
