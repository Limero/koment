package youtube

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/limero/koment/lib/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetch(t *testing.T) {
	// TODO
}

func TestInstances(t *testing.T) {
	if strings.ToLower(os.Getenv("TEST_EXTERNAL")) != "true" {
		t.Skip("Not testing external")
	}

	for _, instance := range defaultInstances {
		t.Run(instance, func(t *testing.T) {
			var resp CommentsResponse
			err := util.GetPageToJSON(fmt.Sprintf(
				"%s/api/v1/comments/dQw4w9WgXcQ/",
				instance,
			), &resp)
			require.NoError(t, err)
			assert.NotEmpty(t, resp.Comments)
			assert.NotEmpty(t, resp.VideoID)
		})
	}
}
