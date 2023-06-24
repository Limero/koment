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
	Use:   "demo",
	Short: "Show fake comments",
	RunE: func(_ *cobra.Command, _ []string) error {
		siteInput := model.SiteInput{
			SiteName: model.SiteDemo,
		}

		app := app.NewApp()
		app.SiteInput = siteInput
		app.RunApp()
		return nil
	},
}
