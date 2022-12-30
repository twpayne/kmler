package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

// A Config is a configuration.
type Config struct {
	output string
	stdout io.WriteCloser
}

type configOption func(*Config)

func newConfig(options ...configOption) *Config {
	c := &Config{
		stdout: os.Stdout,
	}

	for _, option := range options {
		option(c)
	}

	return c
}

func (c *Config) register(rootCmd *cobra.Command) {
	persistentFlags := rootCmd.PersistentFlags()
	persistentFlags.StringVarP(&c.output, "output", "o", c.output, "output file")
	panicOnError(rootCmd.MarkPersistentFlagFilename("output"))
}

func (c *Config) writeOutput(data []byte) error {
	if c.output == "" || c.output == "-" {
		_, err := c.stdout.Write(data)
		return err
	}
	//nolint:gosec
	return os.WriteFile(c.output, data, 0o666)
}

func (c *Config) writeOutputString(s string) error {
	return c.writeOutput([]byte(s))
}
