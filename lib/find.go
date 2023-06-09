package lib

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/limero/koment/lib/model"
	"github.com/limero/koment/lib/site/disqus"
	"github.com/limero/koment/lib/site/reddit"
	"github.com/limero/koment/lib/site/vbulletin"
	"github.com/limero/koment/lib/site/ycombinator"
	"github.com/limero/koment/lib/site/youtube"
)

func FindComments(urlString string) (*model.SiteInput, error) {
	if strings.EqualFold(urlString, "demo") {
		return &model.SiteInput{
			SiteName: model.SiteDemo,
		}, nil
	}

	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	hostname := strings.TrimPrefix(parsedURL.Hostname(), "www.")
	switch hostname {
	case "feber.se", "tjock.se":
		site := disqus.NewDisqus()
		return site.GetInput(parsedURL, "feber")
	case "phoronix.com", "forum.vbulletin.com":
		site := vbulletin.NewVbulletin()
		return site.GetInput(parsedURL)
	case "reddit.com", "old.reddit.com":
		site := reddit.NewReddit()
		return site.GetInput(parsedURL)
	case "news.ycombinator.com":
		site := ycombinator.NewYcombinator()
		return site.GetInput(parsedURL)
	case "youtube.com":
		site := youtube.NewYoutube()
		return site.GetInput(parsedURL)
	}

	return nil, fmt.Errorf("could not find any comments for url %q", urlString)
}
