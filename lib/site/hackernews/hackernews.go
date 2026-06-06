package hackernews

import (
	"fmt"
	"net/url"
	"strconv"
	"sync"

	"github.com/limero/koment/lib/internal/util"
	"github.com/limero/koment/lib/model"
)

type HackerNews struct{}

func NewHackerNews() HackerNews {
	return HackerNews{}
}

func (s HackerNews) GetInput(url *url.URL, _ ...string) (*model.SiteInput, error) {
	queryValues := url.Query()
	return &model.SiteInput{
		SiteName: model.SiteHackerNews,
		ID:       queryValues.Get("id"),
	}, nil
}

func (s HackerNews) Fetch(fi model.SiteInput) (model.Posts, error) {
	return s.getFromAPI(fi)
}

func (s HackerNews) getFromAPI(fi model.SiteInput) (model.Posts, error) {
	id := fi.ID
	depth := 0
	if fi.ContinueFrom != nil {
		id = fi.ContinueFrom.Key
		depth = fi.ContinueFrom.Depth
	}

	apiURL := "https://hacker-news.firebaseio.com/v0/item/%s.json"
	var resp Post
	if err := util.GetPageToJSON(fmt.Sprintf(apiURL, id), &resp); err != nil {
		return nil, err
	}

	var posts model.Posts
	switch resp.Type {
	case "story":
		type result struct {
			pos   int
			posts model.Posts
			err   error
		}
		results := make(chan result, len(resp.Kids))
		sem := make(chan struct{}, 10)
		var wg sync.WaitGroup

		for i, kid := range resp.Kids {
			wg.Add(1)
			go func(pos, kid int) {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()

				var newResp Post
				if err := util.GetPageToJSON(fmt.Sprintf(apiURL, strconv.Itoa(kid)), &newResp); err != nil {
					results <- result{pos: pos, err: err}
					return
				}

				batch := newResp.toModelBatch(depth, resp.By)
				results <- result{pos: pos, posts: batch}
			}(i, kid)
		}

		go func() {
			wg.Wait()
			close(results)
		}()

		pending := make(map[int]model.Posts)
		next := 0
		for res := range results {
			if res.err != nil {
				return nil, res.err
			}
			if res.pos == next {
				posts = append(posts, res.posts...)
				if fi.OnPost != nil {
					fi.OnPost(res.posts)
				}
				next++
				for {
					if p, ok := pending[next]; ok {
						posts = append(posts, p...)
						if fi.OnPost != nil {
							fi.OnPost(p)
						}
						delete(pending, next)
						next++
					} else {
						break
					}
				}
			} else {
				pending[res.pos] = res.posts
			}
		}
	case "comment":
		posts = append(posts, resp.toModelBatch(depth)...)
	}

	return posts, nil
}
