package model

import "net/url"

type SiteName string

const (
	SiteDemo SiteName = "demo"

	SiteDisqus      SiteName = "disqus"
	SiteReddit      SiteName = "reddit"
	SiteYcombinator SiteName = "ycombinator" // Hacker News
	SiteYoutube     SiteName = "youtube"
)

func AllSites() []SiteName {
	return []SiteName{
		SiteDemo,
		SiteDisqus,
		SiteReddit,
		SiteYcombinator,
		SiteYoutube,
	}
}

func AllSitesAsStrings() []string {
	allSites := AllSites()
	sites := make([]string, len(allSites))
	for _, site := range allSites {
		sites = append(sites, string(site))
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
	ContinueFrom *ContinueFrom
	ApiKey       string
}

type ContinueFrom struct {
	Key   string
	Depth int
}
