package wordpress

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/limero/koment/lib/internal/util"
	"github.com/limero/koment/lib/model"
)

type Wordpress struct{}

func NewWordPress() Wordpress {
	return Wordpress{}
}

func (s Wordpress) GetInput(u *url.URL, _ ...string) (*model.SiteInput, error) {
	return &model.SiteInput{
		SiteName: model.SiteWordPress,
		FullURL:  u,
	}, nil
}

func (s Wordpress) Fetch(fi model.SiteInput) (model.Posts, error) {
	html, err := util.GetPageBodyString(fi.FullURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %w", err)
	}

	postID, apiBase, err := extractPostInfo(html)
	if err != nil {
		return nil, fmt.Errorf("failed to extract post info: %w", err)
	}

	apiURL := fmt.Sprintf("%s/comments?post=%s&per_page=100", apiBase, postID)
	var comments []WordPressComment
	if err := util.GetPageToJSON(apiURL, &comments); err != nil {
		return nil, fmt.Errorf("failed to fetch comments: %w", err)
	}

	return toModel(comments), nil
}

func extractPostInfo(html string) (postID string, apiBase string, err error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	sel := doc.Find(`link[rel="alternate"][type="application/json"]`)
	if sel.Length() == 0 {
		return "", "", fmt.Errorf("could not find post ID link tag")
	}

	href, exists := sel.Attr("href")
	if !exists {
		return "", "", fmt.Errorf("link tag missing href attribute")
	}

	u, err := url.Parse(href)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse href URL %q: %w", href, err)
	}

	pathParts := strings.Split(strings.TrimRight(u.Path, "/"), "/")
	if len(pathParts) < 3 {
		return "", "", fmt.Errorf("could not parse post ID from path %q", u.Path)
	}

	postID = pathParts[len(pathParts)-1]
	if postID == "" || pathParts[len(pathParts)-2] != "posts" {
		return "", "", fmt.Errorf("unexpected href format %q", href)
	}

	apiBase = u.Scheme + "://" + u.Host + strings.Join(pathParts[:len(pathParts)-2], "/")
	return postID, apiBase, nil
}
