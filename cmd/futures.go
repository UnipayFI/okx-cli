package cmd

import (
	"github.com/UnipayFI/okx-cli/cmd/futures"
	"github.com/spf13/cobra"
)

// futuresCmd is the parent command for futures trading. The product line is
// selected with the persistent --instType flag (default swap); it applies to
// every subcommand.
var futuresCmd = &cobra.Command{
	Use:   "futures",
	Short: "Futures trading (health, orders & positions)",
	Long: `Futures trading commands for the unified (v5) account: account health,
order entry (create / cancel / get / open), positions (list / close) and the
account position mode (posMode get / set). The product line is selected with
the persistent --instType flag (default swap).
Perpetual-swap instruments are quoted as BASE-QUOTE-SWAP, e.g. BTC-USDT-SWAP.

Docs Link: https://www.okx.com/docs-v5/en/#order-book-trading-trade`,
}

func init() {
	futuresCmd.PersistentFlags().StringP("instType", "C", "swap",
		"instrument type: swap, futures")
	futuresCmd.AddCommand(futures.InitCmds()...)
	RootCmd.AddCommand(futuresCmd)
}
