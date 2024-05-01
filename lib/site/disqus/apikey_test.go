package disqus

import (
	"os"
	"strings"
	"testing"

	"github.com/limero/koment/lib/internal/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_getApiKey(t *testing.T) {
	s := NewDisqus()
	apiKeyFile := helper.CachePath("disqus-api-key.txt")
	defer os.Remove(apiKeyFile)

	t.Run("get cached api key", func(t *testing.T) {
		helper.WriteFile("abc", apiKeyFile)
		apiKey, err := s.getApiKey()
		require.NoError(t, err)
		assert.Equal(t, "abc", apiKey)
	})

	t.Run("fetch api key", func(t *testing.T) {
		if strings.ToLower(os.Getenv("TEST_EXTERNAL")) != "true" {
			t.Skip("Not testing external")
		}
		os.Remove(apiKeyFile)
		apiKey, err := s.getApiKey()
		require.NoError(t, err)
		assert.Len(t, apiKey, 64)
	})
}
