package futures

import (
	"errors"

	"github.com/UnipayFI/okx-cli/config"
	"github.com/UnipayFI/okx-cli/exchange"
	"github.com/UnipayFI/okx-cli/printer"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

const docBase = "https://www.okx.com/docs-v5/en/#order-book-trading-trade-post-"

var (
	healthCmd = &cobra.Command{
		Use:   "health",
		Short: "Show account health (margin ratio & risk)",
		Long: `Show the unified account's risk/health metrics: total and adjusted equity,
unrealised PnL, initial/maintenance margin requirements (IMR/MMR) and the
margin ratio. With --json the position-risk snapshot for the instrument type is
included.

Docs Link: https://www.okx.com/docs-v5/en/#trading-account-rest-api-get-balance`,
		RunE: showHealth,
	}

	orderCmd = &cobra.Command{
		Use:   "order",
		Short: "Create, cancel and query futures orders",
	}

	createCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "Create a futures order",
		Long: `Place a new futures order.
* Required: --instId, --side, --type, --sz
* --px is required for limit orders
* --posSide (long/short) is required in hedge mode
* --tdMode defaults to cross

Docs Link: ` + docBase + "place-order",
		RunE: createOrder,
	}

	cancelCmd = &cobra.Command{
		Use:   "cancel",
		Short: "Cancel a futures order",
		Long: `Cancel a single futures order by --ordId or --clOrdId.

Docs Link: ` + docBase + "cancel-order",
		RunE: cancelOrder,
	}

	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Query a single futures order",
		Long: `Query a single futures order by --ordId or --clOrdId.

Docs Link: https://www.okx.com/docs-v5/en/#order-book-trading-trade-get-order-details`,
		RunE: getOrder,
	}

	openCmd = &cobra.Command{
		Use:   "open",
		Short: "List open futures orders",
		Long: `List currently open (live / partially filled) futures orders.

Docs Link: https://www.okx.com/docs-v5/en/#order-book-trading-trade-get-order-list`,
		RunE: openOrders,
	}

	positionCmd = &cobra.Command{
		Use:   "position",
		Short: "List and close futures positions",
	}

	positionListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List open positions",
		Long: `List the account's open positions for the instrument type.

Docs Link: https://www.okx.com/docs-v5/en/#trading-account-rest-api-get-positions`,
		RunE: listPositions,
	}

	closeCmd = &cobra.Command{
		Use:   "close",
		Short: "Market-close a position",
		Long: `Market-close a position for an instrument under a margin mode.
* Required: --instId, --mgnMode
* --posSide (long/short) is required in hedge mode

Docs Link: ` + docBase + "close-positions",
		RunE: closePosition,
	}
)

// InitCmds registers flags and returns the futures subcommands.
func InitCmds() []*cobra.Command {
	createCmd.Flags().StringP("instId", "s", "", "instrument id, e.g. BTC-USDT-SWAP (required)")
	createCmd.Flags().StringP("side", "S", "", "buy or sell (required)")
	createCmd.Flags().StringP("type", "t", "", "order type: market, limit, post_only, fok, ioc (required)")
	createCmd.Flags().StringP("sz", "q", "", "order size in contracts (decimal) (required)")
	createCmd.Flags().StringP("px", "p", "", "order price (required for limit)")
	createCmd.Flags().StringP("tdMode", "m", "cross", "trade mode: cross, isolated")
	createCmd.Flags().StringP("posSide", "P", "", "position side: long, short (hedge mode)")
	createCmd.Flags().BoolP("reduceOnly", "r", false, "reduce-only order")
	createCmd.Flags().StringP("ccy", "y", "", "margin currency (single-currency margin)")
	createCmd.Flags().String("clOrdId", "", "client order id")
	createCmd.MarkFlagRequired("instId")
	createCmd.MarkFlagRequired("side")
	createCmd.MarkFlagRequired("type")
	createCmd.MarkFlagRequired("sz")

	cancelCmd.Flags().StringP("instId", "s", "", "instrument id, e.g. BTC-USDT-SWAP (required)")
	cancelCmd.Flags().StringP("ordId", "i", "", "order id")
	cancelCmd.Flags().StringP("clOrdId", "c", "", "client order id")
	cancelCmd.MarkFlagRequired("instId")

	getCmd.Flags().StringP("instId", "s", "", "instrument id, e.g. BTC-USDT-SWAP (required)")
	getCmd.Flags().StringP("ordId", "i", "", "order id")
	getCmd.Flags().StringP("clOrdId", "c", "", "client order id")
	getCmd.MarkFlagRequired("instId")

	openCmd.Flags().StringP("instId", "s", "", "instrument id filter, e.g. BTC-USDT-SWAP")

	positionListCmd.Flags().StringP("instId", "s", "", "instrument id filter, e.g. BTC-USDT-SWAP")

	closeCmd.Flags().StringP("instId", "s", "", "instrument id, e.g. BTC-USDT-SWAP (required)")
	closeCmd.Flags().StringP("mgnMode", "m", "cross", "margin mode: cross, isolated")
	closeCmd.Flags().StringP("posSide", "P", "", "position side: long, short (hedge mode)")
	closeCmd.Flags().Bool("autoCxl", false, "cancel pending orders that block the close")
	closeCmd.MarkFlagRequired("instId")

	orderCmd.AddCommand(createCmd, cancelCmd, getCmd, openCmd)
	positionCmd.AddCommand(positionListCmd, closeCmd)
	return []*cobra.Command{healthCmd, orderCmd, positionCmd}
}

// resolveInstType reads the inherited persistent --instType flag (default swap).
func resolveInstType(cmd *cobra.Command) (exchange.InstType, error) {
	raw, _ := cmd.Flags().GetString("instType")
	if raw == "" {
		raw = "swap"
	}
	return exchange.ParseInstType(raw)
}

func showHealth(cmd *cobra.Command, _ []string) error {
	instType, err := resolveInstType(cmd)
	if err != nil {
		return err
	}
	client := exchange.NewClient()
	// In JSON mode surface the richer position-risk snapshot; the table view
	// shows the account-level health summary.
	if config.Config.OutputJSON {
		risk, rerr := client.GetPositionRisk(instType)
		if rerr != nil {
			return rerr
		}
		printer.Print(risk)
		return nil
	}
	bal, err := client.GetBalance()
	if err != nil {
		return err
	}
	view := exchange.HealthView(*bal)
	printer.Print(&view)
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
	p.PosSide, _ = cmd.Flags().GetString("posSide")
	p.ReduceOnly, _ = cmd.Flags().GetBool("reduceOnly")
	p.Ccy, _ = cmd.Flags().GetString("ccy")
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
	instType, err := resolveInstType(cmd)
	if err != nil {
		return err
	}
	instID, _ := cmd.Flags().GetString("instId")
	list, err := exchange.NewClient().GetPendingOrders(instType, instID)
	if err != nil {
		return err
	}
	printer.Print(exchange.OrderRows(list))
	return nil
}

func listPositions(cmd *cobra.Command, _ []string) error {
	instType, err := resolveInstType(cmd)
	if err != nil {
		return err
	}
	instID, _ := cmd.Flags().GetString("instId")
	positions, err := exchange.NewClient().GetPositions(instType, instID)
	if err != nil {
		return err
	}
	printer.Print(exchange.Positions(positions))
	return nil
}

func closePosition(cmd *cobra.Command, _ []string) error {
	instID, _ := cmd.Flags().GetString("instId")
	mgnModeRaw, _ := cmd.Flags().GetString("mgnMode")
	posSide, _ := cmd.Flags().GetString("posSide")
	autoCxl, _ := cmd.Flags().GetBool("autoCxl")
	mgnMode, err := exchange.ParseMgnMode(mgnModeRaw)
	if err != nil {
		return err
	}
	res, err := exchange.NewClient().ClosePosition(instID, mgnMode, posSide, autoCxl)
	if err != nil {
		return err
	}
	view := exchange.ClosePositionView(*res)
	printer.Print(&view)
	return nil
}
