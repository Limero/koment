package reddit

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/limero/koment/lib/internal/util"
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

	fullURL, err := url.Parse(fmt.Sprintf("https://old.reddit.com%s", url.Path))
	if err != nil {
		return nil, fmt.Errorf("failed to parse reddit URL: %w", err)
	}

	return &model.SiteInput{
		SiteName: model.SiteReddit,
		FullURL:  fullURL,
		Category: strings.Split(strings.Split(url.Path, "/r/")[1], "/")[0],
		ID:       strings.Split(strings.Split(url.Path, "/comments/")[1], "/")[0],
	}, nil
}

func (s Reddit) Fetch(fi model.SiteInput) (model.Posts, error) {
	return getFromHTML(fi.FullURL)
}

func getFromHTML(url *url.URL) (model.Posts, error) {
	body, err := util.GetPageBodyString(url.String())
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	return parseComments(doc)
}
