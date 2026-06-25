package cmd

import (
	"github.com/UnipayFI/okx-cli/cmd/spot"
	"github.com/spf13/cobra"
)

// spotCmd is the parent command for spot trading (assets and orders).
var spotCmd = &cobra.Command{
	Use:   "spot",
	Short: "Spot trading (assets & orders)",
	Long: `Spot trading commands for the unified (v5) account: per-currency assets and
order entry (create / cancel / get / open). Spot instruments are quoted as
BASE-QUOTE, e.g. BTC-USDT.

Docs Link: https://www.okx.com/docs-v5/en/#order-book-trading-trade`,
}

func init() {
	spotCmd.AddCommand(spot.InitCmds()...)
	RootCmd.AddCommand(spotCmd)
}
