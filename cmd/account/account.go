package account

import (
	"github.com/UnipayFI/okx-cli/config"
	"github.com/UnipayFI/okx-cli/exchange"
	"github.com/UnipayFI/okx-cli/printer"
	"github.com/spf13/cobra"
)

const docBase = "https://www.okx.com/docs-v5/en/#trading-account-rest-api-"

var (
	infoCmd = &cobra.Command{
		Use:   "info",
		Short: "Show account identity and permissions",
		Long: `Show the unified account's identity, account level, position mode, KYC
level and API permissions.

Docs Link: ` + docBase + "get-account-configuration",
		RunE: showInfo,
	}

	balanceCmd = &cobra.Command{
		Use:   "balance",
		Short: "Show account equity and per-coin balances",
		Long: `Show the unified account's aggregate equity and margin metrics, followed by
the per-currency balance details (non-zero only).

Docs Link: ` + docBase + "get-balance",
		RunE: showBalance,
	}

	healthCmd = &cobra.Command{
		Use:   "health",
		Short: "Show account health (margin ratio & risk)",
		Long: `Show the unified account's risk/health metrics: total and adjusted equity,
unrealised PnL, initial/maintenance margin requirements (IMR/MMR) and the
margin ratio.

Docs Link: ` + docBase + "get-balance",
		RunE: showHealth,
	}
)

// InitCmds registers flags and returns the account subcommands.
func InitCmds() []*cobra.Command {
	return []*cobra.Command{infoCmd, balanceCmd, healthCmd}
}

func showInfo(cmd *cobra.Command, _ []string) error {
	cfg, err := exchange.NewClient().GetAccountConfig()
	if err != nil {
		return err
	}
	view := exchange.AccountConfigView(*cfg)
	printer.Print(&view)
	return nil
}

func showBalance(cmd *cobra.Command, _ []string) error {
	bal, err := exchange.NewClient().GetBalance()
	if err != nil {
		return err
	}
	summary := exchange.BalanceSummary(*bal)
	printer.Print(&summary)
	// In JSON mode the summary already carries the per-coin details; only the
	// table view needs the separate per-coin table.
	if !config.Config.OutputJSON {
		printer.Print(exchange.CoinAssets(bal.Details))
	}
	return nil
}

func showHealth(cmd *cobra.Command, _ []string) error {
	bal, err := exchange.NewClient().GetBalance()
	if err != nil {
		return err
	}
	view := exchange.HealthView(*bal)
	printer.Print(&view)
	return nil
}
