package youtube

import (
	"fmt"
	"net/url"

	"github.com/limero/koment/lib/internal/util"
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
	return s.getFromApi(fi.ID, fi.ContinueFrom)
}

func (s Youtube) getFromApi(videoID string, continueFrom *model.ContinueFrom) (model.Posts, error) {
	continueFromKey := ""
	depth := 0
	if continueFrom != nil {
		continueFromKey = continueFrom.Key
		depth = continueFrom.Depth
	}

	var resp CommentsResponse
	if err := util.GetPageToJSON(fmt.Sprintf(
		"%s/api/v1/comments/%s/?continuation=%s",
		s.invidiousInstance,
		videoID,
		continueFromKey,
	), &resp); err != nil {
		return nil, err
	}

	return resp.toModel(depth)
}
