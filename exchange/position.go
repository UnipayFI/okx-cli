package exchange

import "github.com/UnipayFI/go-okx"

// ---- service calls -------------------------------------------------------

// GetPositions lists the account's open positions, optionally filtered by
// instrument type and instrument id.
func (c *Client) GetPositions(instType InstType, instID string) ([]okx.Position, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.okx.NewGetPositionsService()
	if instType != "" {
		s.SetInstType(instType)
	}
	if instID != "" {
		s.SetInstId(instID)
	}
	return s.Do(cx)
}

// ClosePosition market-closes a position for a symbol under a margin mode. In
// hedge mode posSide selects the side; autoCxl cancels pending orders that would
// otherwise block the close.
func (c *Client) ClosePosition(instID string, mgnMode MgnMode, posSide string, autoCxl bool) (*okx.ClosePositionResult, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.okx.NewClosePositionService(instID, mgnMode)
	if posSide != "" {
		ps, err := ParsePosSide(posSide)
		if err != nil {
			return nil, err
		}
		s.SetPosSide(ps)
	}
	if autoCxl {
		s.SetAutoCxl(true)
	}
	return s.Do(cx)
}

// ---- table models --------------------------------------------------------

// Positions renders the account's open positions.
type Positions []okx.Position

func (p Positions) Header() []string {
	return []string{"Inst ID", "Pos Side", "Margin Mode", "Leverage", "Position", "Available", "Avg Price", "Mark Price", "Liq Price", "Unrealised PNL", "PNL Ratio", "Mgn Ratio", "MMR"}
}

func (p Positions) Row() [][]any {
	rows := [][]any{}
	for _, pos := range p {
		rows = append(rows, []any{
			pos.InstrumentID, pos.PositionSide, pos.MarginMode, pos.Leverage, pos.Position, pos.AvailablePosition,
			pos.AveragePrice, pos.MarkPrice, pos.LiquidationPrice, pos.UPL, pos.UPLRatio, pos.MarginRatio, pos.MMR,
		})
	}
	return rows
}

// ClosePositionView renders the result of a market close.
type ClosePositionView okx.ClosePositionResult

func (c *ClosePositionView) Header() []string {
	return []string{"Inst ID", "Pos Side", "Client Order ID"}
}

func (c *ClosePositionView) Row() [][]any {
	return [][]any{{c.InstrumentID, c.PositionSide, c.ClientOrderID}}
}
