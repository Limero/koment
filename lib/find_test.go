package lib

import (
	"net/url"
	"testing"

	"github.com/limero/koment/lib/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustURL(s string) *url.URL {
	u, _ := url.Parse(s)
	return u
}

func TestFindComments(t *testing.T) {
	for _, tt := range []struct {
		name     string
		url      string
		expected *model.SiteInput
	}{
		{
			name: "demo",
			url:  "demo",
			expected: &model.SiteInput{
				SiteName: model.SiteDemo,
			},
		},
		{
			name: "feber",
			url:  "https://feber.se/abc/def/123456/",
			expected: &model.SiteInput{
				SiteName: model.SiteDisqus,
				ID:       "123456",
				Category: "feber",
			},
		},
		{
			name: "reddit",
			url:  "https://reddit.com/r/subreddit/comments/12dx0b0/abc/",
			expected: &model.SiteInput{
				SiteName: model.SiteReddit,
				Category: "subreddit",
				ID:       "12dx0b0",
				FullURL:  mustURL("https://old.reddit.com/r/subreddit/comments/12dx0b0/abc/"),
			},
		},
		{
			name: "youtube",
			url:  "https://www.youtube.com/watch?v=W0-ql0PiA-U",
			expected: &model.SiteInput{
				SiteName: model.SiteYoutube,
				ID:       "W0-ql0PiA-U",
			},
		},
		{
			name: "liliputing",
			url:  "https://liliputing.com/dell-xps-14-laptop-is-now-available-with-ubuntu-linux/",
			expected: &model.SiteInput{
				SiteName: model.SiteWordPress,
				FullURL:  mustURL("https://liliputing.com/dell-xps-14-laptop-is-now-available-with-ubuntu-linux/"),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := FindComments(tt.url)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
