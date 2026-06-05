package app

import (
	"fmt"

	"github.com/limero/koment/app/ui"
	"github.com/limero/koment/lib"
	"github.com/limero/koment/lib/model"
)

func (a *App) RunHeadless(hu ui.HeadlessUI) error {
	a.Site = lib.NewSite(a.SiteInput.SiteName)

	posts, err := a.Site.Fetch(a.SiteInput)
	if err != nil {
		return fmt.Errorf("error fetching comments: %w", err)
	}

	threads := model.PostsToThreads(posts)
	if len(posts) == 0 {
		fmt.Println("No comments available")
		return nil
	}

	hu.Render(threads)
	return nil
}
