package cmd

import (
	"errors"

	"github.com/UnipayFI/okx-cli/config"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "okx-cli",
	Short: "OKX API CLI for the unified (v5) trading account",
	Long: `okx-cli is a command-line client for OKX's private v5 REST API
(/api/v5/*). OKX runs one unified trading account, so a single client serves
spot and every futures line; commands are grouped by product:

  account  Account identity, balances and health (margin ratio)
  spot     Spot trading: assets, place / cancel / query orders
  futures  Futures trading: health, orders, positions

Credentials are read from the environment:
  OKX_API_KEY, OKX_API_SECRET, OKX_PASSPHRASE   (required)
  OKX_BASE_URL, OKX_DEMO                         (optional)
  HTTPS_PROXY / ALL_PROXY                        (optional proxy)

Use --json on any command for the raw API response.

Docs Link: https://www.okx.com/docs-v5/en/`,
	PersistentPreRunE: checkCredentials,
	SilenceUsage:      true,
	SilenceErrors:     true,
}

func init() {
	initCommandConfig()
}

func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}

func initCommandConfig() {
	RootCmd.CompletionOptions.DisableDefaultCmd = true
	RootCmd.PersistentFlags().BoolVar(&config.Config.OutputJSON, "json", false, "Output JSON instead of a table")
}

func checkCredentials(cmd *cobra.Command, args []string) error {
	if config.Config.APIKey == "" || config.Config.APISecret == "" || config.Config.Passphrase == "" {
		return errors.New("missing credentials: set OKX_API_KEY, OKX_API_SECRET and OKX_PASSPHRASE")
	}
	return nil
}
