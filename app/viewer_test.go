package app

import (
	"testing"

	"github.com/limero/koment/lib/model"
	"github.com/stretchr/testify/assert"
)

func TestContinueStub(t *testing.T) {
	t.Run("Error if demo", func(t *testing.T) {
		a := App{
			Demo: true,
		}
		a.ContinueStub()
		assert.Equal(t, a.infoLevel, "error")
		assert.Contains(t, a.infoMsg, "Fetching replies does not work in demo")
	})

	threads := model.Threads{
		{
			Posts: model.Posts{
				{
					Stub: &model.Stub{},
				},
			},
		},
	}

	t.Run("Error if empty stub key", func(t *testing.T) {
		a := App{
			threads: threads,
		}
		a.ContinueStub()
		assert.Equal(t, a.infoLevel, "error")
		assert.Contains(t, a.infoMsg, "No more replies can be fetched on this thread")
	})
}
