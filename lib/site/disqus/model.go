package disqus

import (
	"time"

	"github.com/limero/koment/lib/model"
)

type Author struct {
	Username                string `json:"username"`
	About                   string `json:"about"`
	Name                    string `json:"name"`
	Disable3RdPartyTrackers bool   `json:"disable3rdPartyTrackers"`
	IsPowerContributor      bool   `json:"isPowerContributor"`
	JoinedAt                string `json:"joinedAt"`
	ProfileURL              string `json:"profileUrl"`
	URL                     string `json:"url"`
	Location                string `json:"location"`
	IsPrivate               bool   `json:"isPrivate"`
	SignedURL               string `json:"signedUrl"`
	IsPrimary               bool   `json:"isPrimary"`
	IsAnonymous             bool   `json:"isAnonymous"`
	ID                      string `json:"id"`
	Avatar                  struct {
		Small struct {
			Permalink string `json:"permalink"`
			Cache     string `json:"cache"`
		} `json:"small"`
		Large struct {
			Permalink string `json:"permalink"`
			Cache     string `json:"cache"`
		} `json:"large"`
		Permalink string `json:"permalink"`
		Cache     string `json:"cache"`
		Xlarge    struct {
			Permalink string `json:"permalink"`
			Cache     string `json:"cache"`
		} `json:"xlarge"`
	} `json:"avatar"`
}

type Post struct {
	EditableUntil          string `json:"editableUntil"`
	Dislikes               int    `json:"dislikes"`
	Thread                 string `json:"thread"`
	NumReports             int    `json:"numReports"`
	Likes                  int    `json:"likes"`
	Message                string `json:"message"`
	ID                     string `json:"id"`
	CreatedAt              string `json:"createdAt"`
	Author                 Author `json:"author"`
	Media                  []any  `json:"media"`
	IsSpam                 bool   `json:"isSpam"`
	IsDeletedByAuthor      bool   `json:"isDeletedByAuthor"`
	IsHighlighted          bool   `json:"isHighlighted"`
	HasMore                bool   `json:"hasMore"`
	Parent                 int64  `json:"parent"`
	IsApproved             bool   `json:"isApproved"`
	IsNewUserNeedsApproval bool   `json:"isNewUserNeedsApproval"`
	IsDeleted              bool   `json:"isDeleted"`
	IsFlagged              bool   `json:"isFlagged"`
	RawMessage             string `json:"raw_message"`
	IsAtFlagLimit          bool   `json:"isAtFlagLimit"`
	CanVote                bool   `json:"canVote"`
	Forum                  string `json:"forum"`
	Depth                  int    `json:"depth"`
	Points                 int    `json:"points"`
	ModerationLabels       []any  `json:"moderationLabels"`
	IsEdited               bool   `json:"isEdited"`
	Sb                     bool   `json:"sb"`
}

type ListPostsThreaded struct {
	Cursor struct {
		Prev    any    `json:"prev"`
		HasNext bool   `json:"hasNext"`
		Next    string `json:"next"`
		HasPrev bool   `json:"hasPrev"`
		Total   int    `json:"total"`
		ID      string `json:"id"`
		More    bool   `json:"more"`
	} `json:"cursor"`
	Code     int    `json:"code"`
	Response []Post `json:"response"`
}

func (from Post) toModel() (model.Post, error) {
	createdAt, err := time.Parse("2006-01-02T15:04:05", from.CreatedAt)
	if err != nil {
		return model.Post{}, err
	}
	return model.Post{
		ID:    from.ID,
		Depth: from.Depth,
		Author: model.Author{
			Name: from.Author.Name,
		},
		Message: from.RawMessage,

		Upvotes:   &from.Likes,
		Downvotes: &from.Dislikes,
		CreatedAt: &createdAt,
	}, nil
}

func (from ListPostsThreaded) toModel() (model.Posts, error) {
	var err error
	posts := make(model.Posts, len(from.Response))
	for i, p := range from.Response {
		posts[i], err = p.toModel()
		if err != nil {
			return nil, err
		}
	}
	return posts, nil
}
