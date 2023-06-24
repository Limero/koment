package demo

import (
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/jaswdr/faker"
	"github.com/limero/koment/lib/model"
)

type Demo struct {
}

func NewDemo() Demo {
	return Demo{}
}

func (s Demo) GetInput(url *url.URL, _ ...string) (*model.SiteInput, error) {
	return &model.SiteInput{
		SiteName: model.SiteDemo,
	}, nil
}

func (s Demo) Fetch(fi model.SiteInput) (model.Posts, error) {
	fake := faker.New()
	posts := make(model.Posts, 0)

	for i := 0; i < fake.IntBetween(10, 20); i++ {
		createdAt := fake.Time().TimeBetween(time.Now().AddDate(0, -1, 0), time.Now())
		posts = append(posts, fakePost(fake, 0, createdAt))

		for j := 0; i < fake.IntBetween(0, 5); j++ {
			createdAt = fake.Time().TimeBetween(createdAt, time.Now())
			posts = append(posts, fakePost(fake, j+1, createdAt))
		}
	}

	return posts, nil
}

func fakePost(fake faker.Faker, depth int, createdAt time.Time) model.Post {
	return model.Post{
		ID:    uuid.NewString(),
		Depth: depth,
		Author: model.Author{
			Name: fake.Person().Name(),
		},
		Message: fake.Lorem().Sentence(fake.IntBetween(3, 50)),

		Upvotes:   nil,
		Downvotes: nil,
		CreatedAt: &createdAt,

		Stub: nil,
	}
}
