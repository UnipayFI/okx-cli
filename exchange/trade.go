package exchange

import (
	"fmt"

	"github.com/UnipayFI/go-okx"
	"github.com/shopspring/decimal"

	"github.com/UnipayFI/okx-cli/common"
)

// ---- service calls -------------------------------------------------------

// PlaceOrderParams collects the fields for a single order. Price is treated as
// unset when zero (required only for limit orders); the futures-only fields are
// applied only when set.
type PlaceOrderParams struct {
	InstID     string
	TdMode     TdMode
	Side       Side
	OrdType    OrdType
	Sz         decimal.Decimal
	Px         decimal.Decimal
	PosSide    string
	ReduceOnly bool
	TgtCcy     string
	Ccy        string
	ClOrdID    string
}

// PlaceOrder submits a new order. OKX returns a per-order sCode even on a
// successful HTTP call, so a non-"0" sCode is surfaced as an error.
func (c *Client) PlaceOrder(p PlaceOrderParams) (*okx.OrderResult, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.okx.NewPlaceOrderService(p.InstID, p.TdMode, p.Side, p.OrdType, p.Sz)
	if !p.Px.IsZero() {
		s.SetPx(p.Px)
	}
	if p.PosSide != "" {
		posSide, err := ParsePosSide(p.PosSide)
		if err != nil {
			return nil, err
		}
		s.SetPosSide(posSide)
	}
	if p.ReduceOnly {
		s.SetReduceOnly(true)
	}
	if p.TgtCcy != "" {
		s.SetTgtCcy(okx.TgtCcy(p.TgtCcy))
	}
	if p.Ccy != "" {
		s.SetCcy(p.Ccy)
	}
	if p.ClOrdID != "" {
		s.SetClOrdId(p.ClOrdID)
	}
	res, err := s.Do(cx)
	if err != nil {
		return nil, err
	}
	if err := checkResult(res.SCode, res.SMsg); err != nil {
		return res, err
	}
	return res, nil
}

// CancelOrder cancels a single order identified by ordId or clOrdId.
func (c *Client) CancelOrder(instID, ordID, clOrdID string) (*okx.OrderResult, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.okx.NewCancelOrderService(instID)
	if ordID != "" {
		s.SetOrdId(ordID)
	}
	if clOrdID != "" {
		s.SetClOrdId(clOrdID)
	}
	res, err := s.Do(cx)
	if err != nil {
		return nil, err
	}
	if err := checkResult(res.SCode, res.SMsg); err != nil {
		return res, err
	}
	return res, nil
}

// GetOrder queries a single order by ordId or clOrdId.
func (c *Client) GetOrder(instID, ordID, clOrdID string) (*okx.Order, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.okx.NewGetOrderService(instID)
	if ordID != "" {
		s.SetOrdId(ordID)
	}
	if clOrdID != "" {
		s.SetClOrdId(clOrdID)
	}
	return s.Do(cx)
}

// GetPendingOrders lists the account's currently open orders, optionally
// filtered by instrument type and instrument id.
func (c *Client) GetPendingOrders(instType InstType, instID string) ([]okx.Order, error) {
	cx, cancel := ctx()
	defer cancel()
	s := c.okx.NewGetOrdersPendingService()
	if instType != "" {
		s.SetInstType(instType)
	}
	if instID != "" {
		s.SetInstId(instID)
	}
	return s.Do(cx)
}

// checkResult turns a non-"0" per-order status code into an error so a failed
// order is never reported as a success.
func checkResult(sCode, sMsg string) error {
	if sCode != "" && sCode != "0" {
		if sMsg == "" {
			return fmt.Errorf("order rejected: sCode=%s", sCode)
		}
		return fmt.Errorf("order rejected: sCode=%s, sMsg=%s", sCode, sMsg)
	}
	return nil
}

// ---- table models --------------------------------------------------------

// OrderResultView renders the identifiers/status returned by place/cancel.
type OrderResultView okx.OrderResult

func (o *OrderResultView) Header() []string {
	return []string{"Order ID", "Client Order ID", "sCode", "sMsg", "Time"}
}

func (o *OrderResultView) Row() [][]any {
	return [][]any{{o.OrderID, o.ClientOrderID, o.SCode, o.SMsg, common.FormatTime(o.Timestamp)}}
}

// OrderRows renders a collection of orders (get / open).
type OrderRows []okx.Order

func (o OrderRows) Header() []string {
	return []string{"Order ID", "Inst ID", "Side", "Pos Side", "Type", "State", "Price", "Size", "Filled", "Avg Price", "Updated"}
}

func (o OrderRows) Row() [][]any {
	rows := [][]any{}
	for _, ord := range o {
		rows = append(rows, []any{
			ord.OrderID, ord.InstrumentID, ord.Side, ord.PositionSide, ord.OrderType, ord.State,
			ord.Price, ord.Size, ord.AccumulatedFillSize, ord.AveragePrice, common.FormatTime(ord.UpdateTime),
		})
	}
	return rows
}
