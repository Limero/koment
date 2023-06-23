package app

import (
	"errors"
	"testing"

	"github.com/limero/koment/app/test"
	"github.com/limero/koment/lib/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestContinueStub(t *testing.T) {
	t.Run("Error if demo", func(t *testing.T) {
		a := App{
			Demo: true,
		}
		a.ContinueStub()
		assert.Equal(t, "error", a.infoLevel)
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
		assert.Equal(t, "error", a.infoLevel)
		assert.Contains(t, a.infoMsg, "No more replies can be fetched on this thread")
	})

	t.Run("Error if failed to fetch", func(t *testing.T) {
		site := new(test.MockSite)

		a := App{
			Site:    site,
			threads: threads,
		}

		threads[0].Posts[0].Stub.Key = "key"
		site.On("Fetch", mock.Anything).Return(model.Posts{}, errors.New("failed to fetch"))
		a.ContinueStub()

		assert.Equal(t, "error", a.infoLevel)
		assert.Contains(t, a.infoMsg, "failed to fetch")
	})

	t.Run("Continue stub", func(t *testing.T) {
		site := new(test.MockSite)

		a := App{
			Site:    site,
			threads: threads,
		}

		threads[0].Posts[0].Stub.Key = "key"
		posts := model.Posts{
			{}, {},
		}
		site.On("Fetch", mock.Anything).Return(posts, nil)
		a.ContinueStub()

		// posts are added and stub is removed
		assert.Len(t, threads[0].Posts, 2)
	})
}
