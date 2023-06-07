package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use: "koment",
	}
)

func Execute() error {
	return rootCmd.Execute()
}
