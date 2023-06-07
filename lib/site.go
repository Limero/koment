package lib

import (
	"github.com/limero/koment/lib/model"
	"github.com/limero/koment/lib/site/disqus"
	"github.com/limero/koment/lib/site/reddit"
	"github.com/limero/koment/lib/site/youtube"
)

func NewSite(siteName model.SiteName) model.Site {
	switch siteName {
	case model.SiteDisqus:
		return disqus.NewDisqus()
	case model.SiteReddit:
		return reddit.NewReddit()
	case model.SiteYoutube:
		return youtube.NewYoutube()
	}

	return nil
}
