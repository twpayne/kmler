package cmd

import "github.com/spf13/cobra"

var customCmd = &cobra.Command{
	Use:   "custom",
	Short: "Generate a custom KML file",
}

func init() {
	rootCmd.AddCommand(customCmd)
}
