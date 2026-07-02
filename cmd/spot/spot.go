package spot

import (
	"errors"

	"github.com/UnipayFI/okx-cli/exchange"
	"github.com/UnipayFI/okx-cli/printer"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

const docBase = "https://www.okx.com/docs-v5/en/#order-book-trading-trade-post-"

// instType is fixed for the spot command group.
const instType = exchange.InstTypeSpot

var (
	assetsCmd = &cobra.Command{
		Use:   "assets",
		Short: "Show spot per-coin balances (non-zero)",
		Long: `Show the unified account's per-currency spot balances. Only coins with a
non-zero equity / cash / available / frozen balance are shown.

Docs Link: https://www.okx.com/docs-v5/en/#trading-account-rest-api-get-balance`,
		RunE: showAssets,
	}

	orderCmd = &cobra.Command{
		Use:   "order",
		Short: "Create, cancel and query spot orders",
	}

	createCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "Create a spot order",
		Long: `Place a new spot order (tdMode cash).
* Required: --instId, --side, --type, --sz
* --px is required for limit orders
* --tgtCcy sets the market-order size unit: base_ccy or quote_ccy

Docs Link: ` + docBase + "place-order",
		RunE: createOrder,
	}

	cancelCmd = &cobra.Command{
		Use:   "cancel",
		Short: "Cancel a spot order",
		Long: `Cancel a single spot order by --ordId or --clOrdId.

Docs Link: ` + docBase + "cancel-order",
		RunE: cancelOrder,
	}

	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Query a single spot order",
		Long: `Query a single spot order by --ordId or --clOrdId.

Docs Link: https://www.okx.com/docs-v5/en/#order-book-trading-trade-get-order-details`,
		RunE: getOrder,
	}

	openCmd = &cobra.Command{
		Use:   "open",
		Short: "List open spot orders",
		Long: `List currently open (live / partially filled) spot orders.

Docs Link: https://www.okx.com/docs-v5/en/#order-book-trading-trade-get-order-list`,
		RunE: openOrders,
	}
)

// InitCmds registers flags and returns the spot subcommands.
func InitCmds() []*cobra.Command {
	createCmd.Flags().StringP("instId", "s", "", "instrument id, e.g. BTC-USDT (required)")
	createCmd.Flags().StringP("side", "S", "", "buy or sell (required)")
	createCmd.Flags().StringP("type", "t", "", "order type: market, limit, post_only, fok, ioc (required)")
	createCmd.Flags().StringP("sz", "q", "", "order size (decimal) (required)")
	createCmd.Flags().StringP("px", "p", "", "order price (required for limit)")
	createCmd.Flags().StringP("tdMode", "m", "cash", "trade mode: cash, cross, isolated, spot_isolated")
	createCmd.Flags().StringP("tgtCcy", "g", "", "market order size unit: base_ccy or quote_ccy")
	createCmd.Flags().String("clOrdId", "", "client order id")
	createCmd.MarkFlagRequired("instId")
	createCmd.MarkFlagRequired("side")
	createCmd.MarkFlagRequired("type")
	createCmd.MarkFlagRequired("sz")

	cancelCmd.Flags().StringP("instId", "s", "", "instrument id, e.g. BTC-USDT (required)")
	cancelCmd.Flags().StringP("ordId", "i", "", "order id")
	cancelCmd.Flags().StringP("clOrdId", "c", "", "client order id")
	cancelCmd.MarkFlagRequired("instId")

	getCmd.Flags().StringP("instId", "s", "", "instrument id, e.g. BTC-USDT (required)")
	getCmd.Flags().StringP("ordId", "i", "", "order id")
	getCmd.Flags().StringP("clOrdId", "c", "", "client order id")
	getCmd.MarkFlagRequired("instId")

	openCmd.Flags().StringP("instId", "s", "", "instrument id filter, e.g. BTC-USDT")

	orderCmd.AddCommand(createCmd, cancelCmd, getCmd, openCmd)
	return []*cobra.Command{assetsCmd, orderCmd}
}

func showAssets(cmd *cobra.Command, _ []string) error {
	bal, err := exchange.NewClient().GetBalance()
	if err != nil {
		return err
	}
	printer.Print(exchange.CoinAssets(bal.Details))
	return nil
}

func createOrder(cmd *cobra.Command, _ []string) error {
	instID, _ := cmd.Flags().GetString("instId")
	sideRaw, _ := cmd.Flags().GetString("side")
	typeRaw, _ := cmd.Flags().GetString("type")
	szRaw, _ := cmd.Flags().GetString("sz")
	tdModeRaw, _ := cmd.Flags().GetString("tdMode")

	side, err := exchange.ParseSide(sideRaw)
	if err != nil {
		return err
	}
	ordType, err := exchange.ParseOrderType(typeRaw)
	if err != nil {
		return err
	}
	tdMode, err := exchange.ParseTdMode(tdModeRaw)
	if err != nil {
		return err
	}
	sz, err := decimal.NewFromString(szRaw)
	if err != nil {
		return errors.New("invalid --sz: " + err.Error())
	}

	p := exchange.PlaceOrderParams{
		InstID:  instID,
		TdMode:  tdMode,
		Side:    side,
		OrdType: ordType,
		Sz:      sz,
	}
	if pxRaw, _ := cmd.Flags().GetString("px"); pxRaw != "" {
		px, perr := decimal.NewFromString(pxRaw)
		if perr != nil {
			return errors.New("invalid --px: " + perr.Error())
		}
		p.Px = px
	}
	p.TgtCcy, _ = cmd.Flags().GetString("tgtCcy")
	p.ClOrdID, _ = cmd.Flags().GetString("clOrdId")

	res, err := exchange.NewClient().PlaceOrder(p)
	if err != nil {
		return err
	}
	view := exchange.OrderResultView(*res)
	printer.Print(&view)
	return nil
}

func cancelOrder(cmd *cobra.Command, _ []string) error {
	instID, _ := cmd.Flags().GetString("instId")
	ordID, _ := cmd.Flags().GetString("ordId")
	clOrdID, _ := cmd.Flags().GetString("clOrdId")
	if ordID == "" && clOrdID == "" {
		return errors.New("one of --ordId or --clOrdId is required")
	}
	res, err := exchange.NewClient().CancelOrder(instID, ordID, clOrdID)
	if err != nil {
		return err
	}
	view := exchange.OrderResultView(*res)
	printer.Print(&view)
	return nil
}

func getOrder(cmd *cobra.Command, _ []string) error {
	instID, _ := cmd.Flags().GetString("instId")
	ordID, _ := cmd.Flags().GetString("ordId")
	clOrdID, _ := cmd.Flags().GetString("clOrdId")
	if ordID == "" && clOrdID == "" {
		return errors.New("one of --ordId or --clOrdId is required")
	}
	order, err := exchange.NewClient().GetOrder(instID, ordID, clOrdID)
	if err != nil {
		return err
	}
	printer.Print(exchange.OrderRows{*order})
	return nil
}

func openOrders(cmd *cobra.Command, _ []string) error {
	instID, _ := cmd.Flags().GetString("instId")
	list, err := exchange.NewClient().GetPendingOrders(instType, instID)
	if err != nil {
		return err
	}
	printer.Print(exchange.OrderRows(list))
	return nil
}
