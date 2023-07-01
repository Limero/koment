package lib

import (
	"testing"

	"github.com/limero/koment/lib/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
		// TODO: Mock Disqus
		/*
			{
				name: "feber",
				url:  "https://feber.se/abc/def/123456/",
				expected: &model.SiteInput{
					SiteName: model.SiteDisqus,
					ID:   "abc",
				},
			},
		*/
		{
			name: "reddit",
			url:  "https://reddit.com/r/subreddit/comments/12dx0b0/abc/",
			expected: &model.SiteInput{
				SiteName: model.SiteReddit,
				Category: "subreddit",
				ID:       "12dx0b0",
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
	} {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := FindComments(tt.url)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
