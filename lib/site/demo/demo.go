package demo

import (
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/limero/koment/lib/model"
)

type Demo struct{}

func NewDemo() Demo {
	return Demo{}
}

func (s Demo) GetInput(url *url.URL, _ ...string) (*model.SiteInput, error) {
	return &model.SiteInput{
		SiteName: model.SiteDemo,
	}, nil
}

func (s Demo) Fetch(_ model.SiteInput) (model.Posts, error) {
	now := time.Now()
	posts := model.Posts{
		{
			ID:        uuid.NewString(),
			Depth:     0,
			Author:    model.Author{Name: "Alice"},
			Message:   "I think this approach has some serious flaws. Has anyone considered the performance implications of doing it this way?",
			CreatedAt: timePtr(now.Add(-2 * time.Hour)),
		},
		{
			ID:        uuid.NewString(),
			Depth:     1,
			Author:    model.Author{Name: "Bob"},
			Message:   "Great point Alice. We actually benchmarked this last quarter and saw a 40% slowdown compared to the alternative.",
			CreatedAt: timePtr(now.Add(-90 * time.Minute)),
		},
		{
			ID:        uuid.NewString(),
			Depth:     2,
			Author:    model.Author{Name: "Carol"},
			Message:   "Do you have the numbers handy? I'd love to see the benchmark results.",
			CreatedAt: timePtr(now.Add(-60 * time.Minute)),
		},
		{
			ID:        uuid.NewString(),
			Depth:     1,
			Author:    model.Author{Name: "Dave"},
			Message:   "I disagree. The convenience outweighs the performance cost in most cases. Premature optimization is the root of all evil.",
			CreatedAt: timePtr(now.Add(-45 * time.Minute)),
		},
		{
			ID:        uuid.NewString(),
			Depth:     2,
			Author:    model.Author{Name: "Eve"},
			Message:   "But this isn't premature -- we already know this is the hot path from profiling.",
			CreatedAt: timePtr(now.Add(-30 * time.Minute)),
		},
		{
			ID:        uuid.NewString(),
			Depth:     0,
			Author:    model.Author{Name: "Frank"},
			Message:   "Has anyone looked at the new library that just dropped? It supposedly handles this pattern much better.",
			CreatedAt: timePtr(now.Add(-20 * time.Minute)),
		},
		{
			ID:        uuid.NewString(),
			Depth:     1,
			Author:    model.Author{Name: "Grace"},
			Message:   "Yes! I tried it yesterday. The API is much cleaner and the defaults are actually sane.",
			CreatedAt: timePtr(now.Add(-10 * time.Minute)),
		},
	}

	return posts, nil
}

// TODO: When bumping to Go 1.26, use new() instead of this
func timePtr(t time.Time) *time.Time {
	return &t
}
