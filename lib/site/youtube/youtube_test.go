package youtube

import (
	"os"
	"testing"

	"github.com/limero/koment/lib/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetch(t *testing.T) {
	os.Chdir("../../../")
	s := NewYoutube()
	posts, err := s.Fetch(model.SiteInput{
		Demo: true,
	})
	require.NoError(t, err)
	assert.Greater(t, len(posts), 1)
}
