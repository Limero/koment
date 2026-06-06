package lobsters

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/limero/koment/lib/internal/util"
	"github.com/limero/koment/lib/model"
)

type Lobsters struct{}

func NewLobsters() Lobsters {
	return Lobsters{}
}

func (s Lobsters) GetInput(url *url.URL, _ ...string) (*model.SiteInput, error) {
	path := strings.Trim(url.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[0] != "s" {
		return nil, fmt.Errorf("could not find short_id in url %q", url.String())
	}

	return &model.SiteInput{
		SiteName: model.SiteLobsters,
		ID:       parts[1],
	}, nil
}

func (s Lobsters) Fetch(fi model.SiteInput) (model.Posts, error) {
	return s.getFromAPI(fi)
}

func (s Lobsters) getFromAPI(fi model.SiteInput) (model.Posts, error) {
	var story Story
	apiURL := fmt.Sprintf("https://lobste.rs/s/%s.json", fi.ID)
	if err := util.GetPageToJSON(apiURL, &story); err != nil {
		return nil, err
	}

	return story.toModel()
}
