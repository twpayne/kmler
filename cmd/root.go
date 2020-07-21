package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:           "kmler",
	Short:         "Convert tracks, routes, and waypoints to KML",
	SilenceErrors: true,
	SilenceUsage:  true,
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
