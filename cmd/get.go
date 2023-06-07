package cmd

import (
	"github.com/limero/koment/app"
	"github.com/limero/koment/lib"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get comments from url and display them",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		siteInput, err := lib.FindComments(args[0])
		if err != nil {
			return err
		}

		app := app.NewApp()
		app.SiteInput = *siteInput
		app.RunApp()
		return nil
	},
}
