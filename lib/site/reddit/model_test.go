package reddit

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

const testHTML = `<html><body><div id="siteTable_t3_abc" class="sitetable nestedlisting">
<div class="thing id-t1_c1 comment" id="thing_t1_c1" data-fullname="t1_c1" data-author="user1">
<div class="entry unvoted">
<p class="tagline">
<span class="score unvoted" title="10">10 points</span>
<time datetime="2026-06-05T07:23:42+00:00" class="live-timestamp">10 hours ago</time>
</p>
<div class="usertext-body may-blank-within md-container">
<div class="md"><p>parent comment body</p></div>
</div>
</div>
<div class="child">
<div id="siteTable_t1_c1" class="sitetable listing">
<div class="thing id-t1_c2 comment" id="thing_t1_c2" data-fullname="t1_c2" data-author="user2">
<div class="entry unvoted">
<p class="tagline">
<span class="score unvoted" title="5">5 points</span>
<time datetime="2026-06-05T10:43:35+00:00" class="live-timestamp">7 hours ago</time>
</p>
<div class="usertext-body may-blank-within md-container">
<div class="md"><p>reply body</p></div>
</div>
</div>
<div class="child"></div>
</div>
</div>
</div>
</div>
<div class="thing id-t1_c3 comment" id="thing_t1_c3" data-fullname="t1_c3" data-author="user3">
<div class="entry unvoted">
<p class="tagline">
<span class="score unvoted" title="3">3 points</span>
<time datetime="2026-06-05T14:14:44+00:00" class="live-timestamp">3 hours ago</time>
</p>
<div class="usertext-body may-blank-within md-container">
<div class="md"><p>second parent comment</p></div>
</div>
</div>
<div class="child"></div>
</div>
</div></body></html>`

func TestParseComments(t *testing.T) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(testHTML))
	assert.NoError(t, err)

	posts, err := parseComments(doc)
	assert.NoError(t, err)
	assert.Len(t, posts, 3)

	t.Run("parent comment", func(t *testing.T) {
		assert.Equal(t, "t1_c1", posts[0].ID)
		assert.Equal(t, 0, posts[0].Depth)
		assert.Equal(t, "user1", posts[0].Author.Name)
		assert.Equal(t, 10, *posts[0].Upvotes)
		assert.NotNil(t, posts[0].CreatedAt)
		assert.Contains(t, posts[0].Message, "parent comment body")
	})

	t.Run("nested reply", func(t *testing.T) {
		assert.Equal(t, "t1_c2", posts[1].ID)
		assert.Equal(t, 1, posts[1].Depth)
		assert.Equal(t, "user2", posts[1].Author.Name)
		assert.Equal(t, 5, *posts[1].Upvotes)
		assert.Contains(t, posts[1].Message, "reply body")
	})

	t.Run("second parent comment", func(t *testing.T) {
		assert.Equal(t, "t1_c3", posts[2].ID)
		assert.Equal(t, 0, posts[2].Depth)
		assert.Equal(t, "user3", posts[2].Author.Name)
		assert.Equal(t, 3, *posts[2].Upvotes)
		assert.Contains(t, posts[2].Message, "second parent comment")
	})
}
