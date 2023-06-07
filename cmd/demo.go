package cmd

import (
	"github.com/limero/koment/app"
	"github.com/limero/koment/lib/model"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(demoCmd)
}

var demoCmd = &cobra.Command{
	Use:       "demo",
	Short:     "Show fake comments",
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: model.AllSitesAsStrings(),
	RunE: func(_ *cobra.Command, args []string) error {
		site := model.SiteDisqus
		if len(args) == 1 {
			site = model.SiteName(args[0])
		}

		siteInput := model.SiteInput{
			SiteName: site,
		}

		app := app.NewApp()
		app.SiteInput = siteInput
		app.Demo = true
		app.RunApp()
		return nil
	},
}
