package model

import "net/url"

type SiteName string

const (
	SiteDemo SiteName = "demo"

	SiteDisqus     SiteName = "disqus"
	SiteReddit     SiteName = "reddit"
	SiteVbulletin  SiteName = "vbulletin"
	SiteWordPress  SiteName = "wordpress"
	SiteHackerNews SiteName = "hackernews"
	SiteYoutube    SiteName = "youtube"
	SiteLobsters   SiteName = "lobsters"
)

func AllSites() []SiteName {
	return []SiteName{
		SiteDemo,
		SiteDisqus,
		SiteReddit,
		SiteVbulletin,
		SiteWordPress,
		SiteHackerNews,
		SiteYoutube,
		SiteLobsters,
	}
}

func AllSitesAsStrings() []string {
	allSites := AllSites()
	sites := make([]string, len(allSites))
	for i, site := range allSites {
		sites[i] = string(site)
	}
	return sites
}

type Site interface {
	GetInput(url *url.URL, v ...string) (*SiteInput, error)
	Fetch(fi SiteInput) (Posts, error)
}

type SiteInput struct {
	SiteName     SiteName
	ID           string
	Category     string
	FullURL      *url.URL
	ContinueFrom *ContinueFrom
	APIKey       string
	OnPost       func(posts Posts)
}

type ContinueFrom struct {
	Key   string
	Depth int
}
