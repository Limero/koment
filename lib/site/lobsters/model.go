package lobsters

import (
	"time"

	"github.com/limero/koment/lib/model"
)

type Story struct {
	ShortID       string    `json:"short_id"`
	Title         string    `json:"title"`
	URL           string    `json:"url"`
	Score         int       `json:"score"`
	SubmitterUser string    `json:"submitter_user"`
	Tags          []string  `json:"tags"`
	Comments      []Comment `json:"comments"`
}

type Comment struct {
	ShortID        string  `json:"short_id"`
	CreatedAt      string  `json:"created_at"`
	LastEditedAt   *string `json:"last_edited_at"`
	Score          int     `json:"score"`
	ParentComment  *string `json:"parent_comment"`
	Depth          int     `json:"depth"`
	Comment        string  `json:"comment"`
	CommentPlain   string  `json:"comment_plain"`
	CommentingUser string  `json:"commenting_user"`
	IsDeleted      bool    `json:"is_deleted"`
	IsModerated    bool    `json:"is_moderated"`
}

func (story Story) toModel() (model.Posts, error) {
	posts := make(model.Posts, len(story.Comments))
	for i, c := range story.Comments {
		post, err := c.toModel(story.SubmitterUser)
		if err != nil {
			return nil, err
		}
		posts[i] = post
	}
	return posts, nil
}

func (c Comment) toModel(opName string) (model.Post, error) {
	createdAt, err := time.Parse(time.RFC3339Nano, c.CreatedAt)
	if err != nil {
		return model.Post{}, err
	}

	message := c.CommentPlain
	if message == "" {
		message = c.Comment
	}

	return model.Post{
		ID:     c.ShortID,
		Depth:  c.Depth,
		Hidden: c.IsDeleted || c.IsModerated,
		Author: model.Author{
			Name: c.CommentingUser,
		},
		Message:   message,
		IsOP:      c.CommentingUser == opName,
		Upvotes:   &c.Score,
		CreatedAt: &createdAt,
	}, nil
}
