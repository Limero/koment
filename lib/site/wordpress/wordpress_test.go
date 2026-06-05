package wordpress

import (
	"os"
	"testing"

	"github.com/limero/koment/lib/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetch(t *testing.T) {
	if os.Getenv("TEST_EXTERNAL") == "" {
		t.Skip("Not testing external")
	}

	site := NewWordPress()
	input := &model.SiteInput{
		FullUrl: mustURL("https://liliputing.com/dell-xps-14-laptop-is-now-available-with-ubuntu-linux/"),
	}

	posts, err := site.Fetch(*input)
	require.NoError(t, err)
	require.NotEmpty(t, posts)

	for _, p := range posts {
		assert.NotEmpty(t, p.ID, "post ID should not be empty")
		assert.NotEmpty(t, p.Message, "message should not be empty")
		assert.GreaterOrEqual(t, p.Depth, 0, "depth should be >= 0")
	}
}
