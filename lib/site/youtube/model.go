package youtube

import (
	"time"

	"github.com/google/uuid"
	"github.com/limero/koment/lib/model"
)

type Comment struct {
	Verified         bool   `json:"verified"`
	Author           string `json:"author"`
	AuthorThumbnails []struct {
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"authorThumbnails"`
	AuthorID             string `json:"authorId"`
	AuthorURL            string `json:"authorUrl"`
	IsEdited             bool   `json:"isEdited"`
	Content              string `json:"content"`
	ContentHTML          string `json:"contentHtml"`
	IsPinned             bool   `json:"isPinned"`
	IsSponsor            bool   `json:"isSponsor"`
	Published            int64  `json:"published"`
	PublishedText        string `json:"publishedText"`
	LikeCount            int    `json:"likeCount"`
	CommentID            string `json:"commentId"`
	AuthorIsChannelOwner bool   `json:"authorIsChannelOwner"`
	CreatorHeart         struct {
		CreatorThumbnail string `json:"creatorThumbnail"`
		CreatorName      string `json:"creatorName"`
	} `json:"creatorHeart,omitempty"`
	Replies struct {
		ReplyCount   int    `json:"replyCount"`
		Continuation string `json:"continuation"`
	} `json:"replies,omitempty"`
}

type CommentsResponse struct {
	CommentCount int       `json:"commentCount"`
	VideoID      string    `json:"videoId"`
	Comments     []Comment `json:"comments"`
	Continuation string    `json:"continuation"`
}

func (from Comment) toModel(depth int) (model.Post, error) {
	createdAt := time.Unix(from.Published, 0)

	return model.Post{
		ID:    from.CommentID,
		Depth: depth,
		Author: model.Author{
			Name: from.Author,
		},
		Message: from.Content,

		Upvotes:   &from.LikeCount,
		Downvotes: nil, // Dislike count is not public
		CreatedAt: &createdAt,
	}, nil
}

func (from CommentsResponse) toModel(depth int) (model.Posts, error) {
	posts := make(model.Posts, 0)
	for _, p := range from.Comments {
		post, err := p.toModel(depth)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)

		if p.Replies.ReplyCount > 0 {
			posts = append(posts, model.Post{
				ID: uuid.NewString(),
				Stub: &model.Stub{
					Count: p.Replies.ReplyCount,
					Key:   p.Replies.Continuation,
				},
			})
		}
	}
	return posts, nil
}
