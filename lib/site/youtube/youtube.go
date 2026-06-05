package youtube

import (
	"fmt"
	"net/url"
	"os"

	"github.com/limero/koment/lib/internal/util"
	"github.com/limero/koment/lib/model"
)

var defaultInstances = []string{
	"https://inv.nadeko.net",
	"https://yt.chocolatemoo53.com",
}

type Youtube struct {
	instances []string
}

func NewYoutube() Youtube {
	instances := defaultInstances
	if env := os.Getenv("INVIDIOUS_INSTANCE"); env != "" {
		instances = []string{env}
	}
	return Youtube{
		instances: instances,
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
	return s.getFromApi(fi.ID, fi.ContinueFrom)
}

func (s Youtube) getFromApi(videoID string, continueFrom *model.ContinueFrom) (model.Posts, error) {
	continueFromKey := ""
	depth := 0
	if continueFrom != nil {
		continueFromKey = continueFrom.Key
		depth = continueFrom.Depth
	}

	var errs []error
	for _, instance := range s.instances {
		var resp CommentsResponse
		err := util.GetPageToJSON(fmt.Sprintf(
			"%s/api/v1/comments/%s/?continuation=%s",
			instance,
			videoID,
			continueFromKey,
		), &resp)
		if err == nil {
			return resp.toModel(depth)
		}
		errs = append(errs, fmt.Errorf("%s: %w", instance, err))
	}

	return nil, fmt.Errorf("all invidious instances failed:\n  %s", func() string {
		var s string
		for _, e := range errs {
			s += fmt.Sprintf("  - %v\n", e)
		}
		return s
	}())
}
