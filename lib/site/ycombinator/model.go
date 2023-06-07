package ycombinator

import (
	"strconv"
	"time"

	"github.com/limero/koment/lib/model"
)

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

func (from Post) toModel() (model.Post, error) {
	createdAt := time.Unix(from.Time, 0)

	return model.Post{
		ID:    strconv.Itoa(from.ID),
		Depth: 0, // TODO
		Author: model.Author{
			Name: from.By,
		},
		Message: from.Text, // TODO: Cleanup HTML from it

		Upvotes:   &from.Score,
		CreatedAt: &createdAt,
	}, nil
}
