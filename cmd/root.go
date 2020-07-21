package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:           "kmler",
	Short:         "Convert tracks, routes, and waypoints to KML",
	SilenceErrors: true,
	SilenceUsage:  true,
}

var config = newConfig()

func init() {
	config.register(rootCmd)
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
