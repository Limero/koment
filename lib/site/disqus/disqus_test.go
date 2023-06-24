package disqus

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetch(t *testing.T) {
	// TODO
}

func TestGetThreadIDFromEmbedPage(t *testing.T) {
	if strings.ToLower(os.Getenv("TEST_EXTERNAL")) != "true" {
		t.Skip("Not testing external")
	}

	s := NewDisqus()
	threadID, err := s.getThreadIDFromEmbedPage("feber", "450517")
	require.NoError(t, err)
	assert.Equal(t, "9676608399", threadID)
}
