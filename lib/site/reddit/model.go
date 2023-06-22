package reddit

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/limero/koment/lib/model"
)

type Listings []Listing

type Listing struct {
	Data ListingData `json:"data"`
}

type ListingData struct {
	Children []Children `json:"children"`
}

type Children struct {
	Data ChildrenData `json:"data"`
}

type ChildrenData struct {
	ID     string `json:"id"`
	Body   string `json:"body"`
	Depth  int    `json:"depth"`
	Author string `json:"author"`
	// created might have one decimal 0, so can't use int64 directly
	CreatedUTC json.Number `json:"created_utc"`
	Ups        int         `json:"ups"`
	Downs      int         `json:"downs"`

	// replies is sometimes an object and sometimes an empty string
	RepliesRaw json.RawMessage `json:"replies"`
}

func (from Children) toModel() (model.Post, error) {
	createdAt, err := from.getCreatedAt()
	if err != nil {
		return model.Post{}, err
	}

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

func (from Children) getCreatedAt() (time.Time, error) {
	createdAtString, _ := strings.CutSuffix(string(from.Data.CreatedUTC), ".0")
	createdAtInt, err := strconv.ParseInt(createdAtString, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(createdAtInt, 0), nil
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
