package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var routeCmd = &cobra.Command{
	Use:   "route",
	Args:  cobra.ExactArgs(1),
	Short: "Convert a route to KML",
	RunE:  config.runRouteCmdE,
}

func init() {
	rootCmd.AddCommand(routeCmd)
}

func (c *Config) runRouteCmdE(cmd *cobra.Command, args []string) error {
	switch strings.ToLower(filepath.Ext(args[0])) {
	case ".xctsk":
		return c.makeXCTrackTaskRoute(args[0])
	default:
		return fmt.Errorf("%s: unsupported extension", args[1])
	}
}
