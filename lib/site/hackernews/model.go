package hackernews

import (
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/limero/koment/lib/internal/util"
	"github.com/limero/koment/lib/model"
)

type Posts []Post

type Post struct {
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	ID          int    `json:"id"`
	Kids        []int  `json:"kids"`
	Parent      int    `json:"parent"`
	Text        string `json:"text"`
	Score       int    `json:"score"`
	Time        int64  `json:"time"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	URL         string `json:"url"`
}

func (from Post) toModelBatch(depth int, opName ...string) model.Posts {
	post, _ := from.toModel(depth, opName...)
	posts := model.Posts{post}
	for _, kid := range from.Kids {
		posts = append(posts, model.Post{
			ID:    uuid.NewString(),
			Depth: depth + 1,
			Stub: &model.Stub{
				Count: 1,
				Key:   strconv.Itoa(kid),
			},
		})
	}
	return posts
}

func (from Post) toModel(depth int, opName ...string) (model.Post, error) {
	createdAt := time.Unix(from.Time, 0)

	message := util.CleanHTML(from.Text)

	isOP := len(opName) > 0 && opName[0] == from.By

	return model.Post{
		ID:    strconv.Itoa(from.ID),
		Depth: depth,
		Author: model.Author{
			Name: from.By,
		},
		Message: message,
		IsOP:    isOP,

		Upvotes:   &from.Score,
		CreatedAt: &createdAt,
	}, nil
}

func (from Posts) toModel(depth int, opName ...string) (model.Posts, error) {
	posts := make(model.Posts, 0)
	for _, p := range from {
		post, err := p.toModel(depth, opName...)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)

		if len(p.Kids) > 0 {
			for _, kid := range p.Kids {
				posts = append(posts, model.Post{
					ID:    uuid.NewString(),
					Depth: depth + 1,
					Stub: &model.Stub{
						Count: 1,
						Key:   strconv.Itoa(kid),
					},
				})
			}
		}
	}
	return posts, nil
}
