package reddit

import (
	"encoding/json"
	"time"

	"github.com/limero/koment/lib/model"
)

type Listings []Listing

type Listing struct {
	Data struct {
		Children []Children `json:"children"`
	} `json:"data"`
}

type Children struct {
	Data struct {
		ID         string `json:"id"`
		Body       string `json:"body"`
		Depth      int    `json:"depth"`
		Author     string `json:"author"`
		CreatedUTC int64  `json:"created_utc"`
		Ups        int    `json:"ups"`
		Downs      int    `json:"downs"`

		// replies is sometimes an object and sometimes an empty string
		RepliesRaw json.RawMessage `json:"replies"`
	} `json:"data"`
}

func (from Children) toModel() (model.Post, error) {
	createdAt := time.Unix(from.Data.CreatedUTC, 0)

	return model.Post{
		ID:    from.Data.ID,
		Depth: from.Data.Depth,
		Author: model.Author{
			Name: from.Data.Author,
		},
		Message: from.Data.Body,

		Upvotes:   &from.Data.Ups,
		Downvotes: &from.Data.Downs,
		CreatedAt: &createdAt,
	}, nil
}

func (from Listings) toModel() (model.Posts, error) {
	var posts model.Posts
	// First index is just metadata on post
	for _, p := range from[1].Data.Children {
		post, err := p.toModel()
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)

		var replies Listing
		if err := json.Unmarshal(p.Data.RepliesRaw, &replies); err != nil {
			continue
		}
		for _, reply := range replies.Data.Children {
			post, err := reply.toModel()
			if err != nil {
				return nil, err
			}
			posts = append(posts, post)
		}
	}
	return posts, nil
}
