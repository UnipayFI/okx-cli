package cmd

import (
	"github.com/UnipayFI/okx-cli/cmd/account"
	"github.com/spf13/cobra"
)

// accountCmd is the parent command for unified-account identity, balance and
// health (margin-ratio) queries.
var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Account identity, balances and health",
	Long: `Account commands for the unified (v5) trading account: identity and
permissions, aggregate balance/equity, and account health (margin ratio).

Docs Link: https://www.okx.com/docs-v5/en/#trading-account-rest-api`,
}

func init() {
	accountCmd.AddCommand(account.InitCmds()...)
	RootCmd.AddCommand(accountCmd)
}
