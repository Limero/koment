package wordpress

import (
	"sort"
	"strconv"
	"time"

	"github.com/limero/koment/lib/internal/util"
	"github.com/limero/koment/lib/model"
)

type WordPressComment struct {
	ID         int    `json:"id"`
	Post       int    `json:"post"`
	Parent     int    `json:"parent"`
	AuthorName string `json:"author_name"`
	DateGmt    string `json:"date_gmt"`
	Content    struct {
		Rendered string `json:"rendered"`
	} `json:"content"`
}

func toModel(comments []WordPressComment) model.Posts {
	byID := make(map[int]WordPressComment, len(comments))
	for _, c := range comments {
		byID[c.ID] = c
	}

	children := make(map[int][]int)
	for _, c := range comments {
		if c.Parent != 0 {
			if _, ok := byID[c.Parent]; ok {
				children[c.Parent] = append(children[c.Parent], c.ID)
			}
		}
	}

	roots := make([]WordPressComment, 0)
	for _, c := range comments {
		if c.Parent == 0 {
			roots = append(roots, c)
		} else if _, ok := byID[c.Parent]; !ok {
			roots = append(roots, c)
		}
	}
	sort.Slice(roots, func(i, j int) bool {
		return roots[i].ID < roots[j].ID
	})

	posts := make(model.Posts, 0, len(comments))
	var dfs func(id int, depth int)
	dfs = func(id int, depth int) {
		c := byID[id]
		createdAt, _ := time.Parse("2006-01-02T15:04:05", c.DateGmt)
		posts = append(posts, model.Post{
			ID:    strconv.Itoa(c.ID),
			Depth: depth,
			Author: model.Author{
				Name: c.AuthorName,
			},
			Message:   util.CleanHTML(c.Content.Rendered),
			CreatedAt: &createdAt,
		})
		childIDs := children[id]
		sort.Slice(childIDs, func(i, j int) bool {
			return childIDs[i] < childIDs[j]
		})
		for _, cid := range childIDs {
			dfs(cid, depth+1)
		}
	}
	for _, root := range roots {
		dfs(root.ID, 0)
	}
	return posts
}
