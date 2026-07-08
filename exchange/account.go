package exchange

import "github.com/UnipayFI/go-okx"

// ---- service calls -------------------------------------------------------

// GetBalance returns the unified account's aggregate equity, margin metrics and
// per-currency balance details.
func (c *Client) GetBalance() (*okx.Balance, error) {
	cx, cancel := ctx()
	defer cancel()
	return c.okx.NewGetBalanceService().Do(cx)
}

// GetAccountConfig returns the account's identity, level, position mode and
// permission metadata.
func (c *Client) GetAccountConfig() (*okx.AccountConfig, error) {
	cx, cancel := ctx()
	defer cancel()
	return c.okx.NewGetAccountConfigService().Do(cx)
}

// SetPositionMode sets the account-wide position mode for derivatives:
// long/short (hedge) or net (one-way). OKX rejects the switch while the account
// holds open positions or pending orders.
func (c *Client) SetPositionMode(mode PosMode) (*okx.PositionMode, error) {
	cx, cancel := ctx()
	defer cancel()
	return c.okx.NewSetPositionModeService(mode).Do(cx)
}

// GetPositionRisk returns a risk snapshot (balances + positions) for an
// instrument type. An empty instType returns risk across all products.
func (c *Client) GetPositionRisk(instType InstType) (*okx.AccountPositionRisk, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.okx.NewGetAccountPositionRiskService()
	if instType != "" {
		s.SetInstType(instType)
	}
	return s.Do(cx)
}

// ---- table models --------------------------------------------------------

// AccountConfigView renders account identity, level and permission metadata.
type AccountConfigView okx.AccountConfig

func (a *AccountConfigView) Header() []string {
	return []string{"UID", "Account Level", "Position Mode", "KYC Level", "Permissions"}
}

func (a *AccountConfigView) Row() [][]any {
	return [][]any{{a.UID, a.AccountLevel, a.PositionMode, a.KYCLevel, a.Perm}}
}

// PositionModeView renders the account-wide position mode.
type PositionModeView okx.PositionMode

func (p *PositionModeView) Header() []string {
	return []string{"Position Mode"}
}

func (p *PositionModeView) Row() [][]any {
	return [][]any{{p.PositionMode}}
}

// BalanceSummary renders the unified account's aggregate equity and margin.
type BalanceSummary okx.Balance

func (b *BalanceSummary) Header() []string {
	return []string{"Total Equity", "Adj Equity", "Avail Equity", "Order Frozen", "IMR", "MMR", "Mgn Ratio", "Unrealised PNL"}
}

func (b *BalanceSummary) Row() [][]any {
	return [][]any{{b.TotalEquity, b.AdjustedEquity, b.AvailableEquity, b.OrderFrozen, b.IMR, b.MMR, b.MarginRatio, b.UPL}}
}

// CoinAssets renders the per-currency balance details (non-zero only).
type CoinAssets []okx.BalanceDetail

func (a CoinAssets) Header() []string {
	return []string{"Coin", "Equity", "Cash Balance", "Available", "Frozen", "Unrealised PNL", "USD Value"}
}

func (a CoinAssets) Row() [][]any {
	rows := [][]any{}
	for _, d := range a {
		if d.Equity.IsZero() && d.CashBalance.IsZero() && d.AvailableBalance.IsZero() && d.FrozenBalance.IsZero() {
			continue
		}
		rows = append(rows, []any{d.Currency, d.Equity, d.CashBalance, d.AvailableBalance, d.FrozenBalance, d.UPL, d.EquityUSD})
	}
	return rows
}

// HealthView renders the unified account's risk/health metrics: total and
// adjusted equity, unrealised PnL, initial/maintenance margin requirements
// (IMR/MMR) and the margin ratio. On OKX a lower margin ratio signals higher
// liquidation risk.
type HealthView okx.Balance

func (h *HealthView) Header() []string {
	return []string{"Total Equity", "Adj Equity", "Avail Equity", "Unrealised PNL", "IMR", "MMR", "Mgn Ratio"}
}

func (h *HealthView) Row() [][]any {
	return [][]any{{h.TotalEquity, h.AdjustedEquity, h.AvailableEquity, h.UPL, h.IMR, h.MMR, h.MarginRatio}}
}
