package cmd

import (
	"github.com/UnipayFI/okx-cli/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	// Override the root credential check: version must work without API keys.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error { return nil },
	Run: func(cmd *cobra.Command, args []string) {
		version.Version()
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
