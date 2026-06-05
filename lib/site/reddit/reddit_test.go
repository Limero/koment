package reddit

import (
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetch(t *testing.T) {
	if os.Getenv("TEST_EXTERNAL") == "" {
		t.Skip("Skipping external test; set TEST_EXTERNAL=true to run")
	}

	fullUrl, err := url.Parse("https://old.reddit.com/r/Music/comments/56cdgm/ama_im_really_rick_astley_i_swear_and_to/")
	require.NoError(t, err)

	posts, err := getFromHTML(fullUrl)
	require.NoError(t, err)
	assert.NotEmpty(t, posts)

	for _, p := range posts {
		assert.NotEmpty(t, p.ID)
		assert.NotEmpty(t, p.Author.Name)
		assert.NotEmpty(t, p.Message)
		assert.NotNil(t, p.Upvotes)
		assert.NotNil(t, p.CreatedAt)
		assert.GreaterOrEqual(t, p.Depth, 0)
	}
}
