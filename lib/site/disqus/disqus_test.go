package disqus

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchEmbedPage(t *testing.T) {
	if os.Getenv("TEST_EXTERNAL") == "" {
		t.Skip("Not testing external")
	}

	s := NewDisqus()
	embed, err := s.fetchEmbedPage("feber", "450517")
	require.NoError(t, err)
	assert.Equal(t, "9676608399", embed.Response.Thread.ID)
	assert.NotEmpty(t, embed.Response.Posts)
}
