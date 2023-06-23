package app

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/limero/koment/app/test"
	"github.com/limero/koment/lib/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestContinueStub_demo(t *testing.T) {
	a := App{
		Demo: true,
	}
	a.ContinueStub()
	assert.Equal(t, "error", a.infoLevel)
	assert.Contains(t, a.infoMsg, "Fetching replies does not work in demo")
}

func TestContinueStub(t *testing.T) {
	for _, tt := range []struct {
		name          string
		threads       model.Threads
		siteFetch     *test.MockedCall[any, model.Posts]
		expectedErr   string
		expectedPosts []string
	}{
		{
			name: "Error if empty stub key",
			threads: model.Threads{
				{
					Posts: model.Posts{
						{
							Stub: &model.Stub{},
						},
					},
				},
			},
			expectedErr: "No more replies can be fetched on this thread",
		},
		{
			name: "Error if failed to fetch",
			threads: model.Threads{
				{
					Posts: model.Posts{
						{
							Stub: &model.Stub{
								Key: "key",
							},
						},
					},
				},
			},
			siteFetch: &test.MockedCall[any, model.Posts]{
				Error: errors.New("failed to fetch"),
			},
			expectedErr: "failed to fetch",
		},
		{
			name: "Continue stub",
			threads: model.Threads{
				{
					Posts: model.Posts{
						{
							ID: uuid.NewString(),
							Stub: &model.Stub{
								Key: "key",
							},
						},
					},
				},
			},
			siteFetch: &test.MockedCall[any, model.Posts]{
				Return: model.Posts{{ID: "1"}, {ID: "2"}},
			},
			expectedPosts: []string{"1", "2"},
		},
		{
			name: "Continue stub with more posts remaining",
			threads: model.Threads{
				{
					Posts: model.Posts{
						{
							ID: uuid.NewString(),
							Stub: &model.Stub{
								Count: 5,
								Key:   "key",
							},
						},
					},
				},
			},
			siteFetch: &test.MockedCall[any, model.Posts]{
				Return: model.Posts{{ID: "1"}, {ID: "2"}},
			},
			expectedPosts: []string{"1", "2", "stub"},
		},
		{
			name: "Continue stub with additional stubs",
			threads: model.Threads{
				{
					Posts: model.Posts{
						{
							ID: "abc", // not uuid for assert to confirm correct stub
							Stub: &model.Stub{
								Key: "key1",
							},
						},
						{
							ID: uuid.NewString(),
							Stub: &model.Stub{
								Key: "key2",
							},
						},
					},
				},
			},
			siteFetch: &test.MockedCall[any, model.Posts]{
				Return: model.Posts{{ID: "1"}, {ID: "2"}},
			},
			expectedPosts: []string{"1", "2", "stub"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			site := new(test.MockSite)

			if tt.siteFetch != nil {
				site.On("Fetch", mock.Anything).
					Return(tt.siteFetch.Return, tt.siteFetch.Error)
			}

			a := App{
				Site:    site,
				threads: tt.threads,
			}

			a.ContinueStub()

			if tt.expectedErr != "" {
				assert.Equal(t, "error", a.infoLevel)
				assert.Contains(t, a.infoMsg, tt.expectedErr)
			} else {
				// posts are added and stub is removed
				assert.Len(t, a.threads[a.activeThread].Posts, len(tt.expectedPosts))
				for i, postID := range tt.expectedPosts {
					// check if correct posts in the right order
					actual := a.threads[a.activeThread].Posts[i]
					if postID == "stub" {
						assert.Len(t, actual.ID, 36) // generated uuid
						assert.NotNil(t, actual.Stub)
					} else {
						assert.Equal(t, postID, actual.ID)
						assert.Nil(t, actual.Stub)
					}
				}
			}
		})
	}
}
