// Package exchange wraps the go-okx SDK with the thin, CLI-friendly surface the
// cmd layer consumes: a single authenticated client constructor, request
// helpers, and TableWriter response models. OKX runs one unified (v5) account,
// so a single client serves spot and every futures line; the product is
// selected per call via the instrument type (instType) and trade mode (tdMode).
package exchange

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/UnipayFI/go-okx"
	"github.com/UnipayFI/go-okx/client"

	"github.com/UnipayFI/okx-cli/config"
)

// requestTimeout bounds every REST call so a stuck connection fails the command
// instead of hanging the terminal.
const requestTimeout = 30 * time.Second

// Re-export the SDK enums so the cmd layer need not import the SDK package.
type (
	InstType = okx.InstType
	Side     = okx.Side
	OrdType  = okx.OrdType
	TdMode   = okx.TdMode
	PosSide  = okx.PosSide
	MgnMode  = okx.MgnMode
	TgtCcy   = okx.TgtCcy
	OrdState = okx.OrdState
)

const (
	InstTypeSpot    = okx.InstTypeSpot
	InstTypeSwap    = okx.InstTypeSwap
	InstTypeFutures = okx.InstTypeFutures
	InstTypeMargin  = okx.InstTypeMargin
)

// Client is the authenticated wrapper around the SDK's unified (v5) client.
type Client struct {
	okx *okx.Client
}

// BuildOptions assembles the shared SDK client options (auth, proxy, base URL,
// demo, silent logging) from the global config.
func BuildOptions() []client.Options {
	opts := []client.Options{
		client.WithAuth(config.Config.APIKey, config.Config.APISecret, config.Config.Passphrase),
		client.WithLogger(silentLogger{}),
	}
	if config.Config.Proxy != "" {
		opts = append(opts, client.WithProxy(config.Config.Proxy))
	}
	if config.Config.BaseURL != "" {
		opts = append(opts, client.WithBaseURL(config.Config.BaseURL))
	}
	if config.Config.Demo {
		opts = append(opts, client.WithDemoTrading(true))
	}
	return opts
}

// NewClient builds an authenticated unified-account client from the global
// config and best-effort syncs the server clock so signed requests carry an
// accepted timestamp (OKX rejects requests skewed more than 30s).
func NewClient() *Client {
	c := okx.NewClient(BuildOptions()...)
	cx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	_ = c.SyncServerTime(cx)
	return &Client{okx: c}
}

// OKX exposes the underlying SDK client for direct service calls.
func (c *Client) OKX() *okx.Client { return c.okx }

// ctx returns a request-scoped context with the standard timeout. Callers must
// defer the returned cancel func.
func ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), requestTimeout)
}

// ParseInstType resolves an instrument-type string (with friendly aliases) to
// the SDK enum (case-insensitive).
func ParseInstType(s string) (okx.InstType, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "spot":
		return okx.InstTypeSpot, nil
	case "swap", "perp", "perpetual":
		return okx.InstTypeSwap, nil
	case "futures", "future":
		return okx.InstTypeFutures, nil
	case "margin":
		return okx.InstTypeMargin, nil
	default:
		return "", fmt.Errorf("invalid instType %q: want one of spot, swap, futures, margin", s)
	}
}

// ParseSide resolves an order side string to the SDK enum (case-insensitive).
func ParseSide(s string) (okx.Side, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "buy":
		return okx.SideBuy, nil
	case "sell":
		return okx.SideSell, nil
	default:
		return "", fmt.Errorf("invalid side %q: want buy or sell", s)
	}
}

// ParseOrderType resolves an order-type string to the SDK enum (case-insensitive).
func ParseOrderType(s string) (okx.OrdType, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "market":
		return okx.OrdTypeMarket, nil
	case "limit":
		return okx.OrdTypeLimit, nil
	case "post_only", "postonly":
		return okx.OrdTypePostOnly, nil
	case "fok":
		return okx.OrdTypeFOK, nil
	case "ioc":
		return okx.OrdTypeIOC, nil
	default:
		return "", fmt.Errorf("invalid order type %q: want one of market, limit, post_only, fok, ioc", s)
	}
}

// ParseTdMode resolves a trade-mode string to the SDK enum (case-insensitive).
func ParseTdMode(s string) (okx.TdMode, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "cash":
		return okx.TdModeCash, nil
	case "cross", "crossed":
		return okx.TdModeCross, nil
	case "isolated", "isolate":
		return okx.TdModeIsolated, nil
	case "spot_isolated":
		return okx.TdModeSpotIsolated, nil
	default:
		return "", fmt.Errorf("invalid tdMode %q: want one of cash, cross, isolated, spot_isolated", s)
	}
}

// ParsePosSide resolves a position-side string to the SDK enum (case-insensitive).
func ParsePosSide(s string) (okx.PosSide, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "long":
		return okx.PosSideLong, nil
	case "short":
		return okx.PosSideShort, nil
	case "net":
		return okx.PosSideNet, nil
	default:
		return "", fmt.Errorf("invalid posSide %q: want one of long, short, net", s)
	}
}

// ParseMgnMode resolves a margin-mode string to the SDK enum (case-insensitive).
func ParseMgnMode(s string) (okx.MgnMode, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "cross", "crossed":
		return okx.MgnModeCross, nil
	case "isolated", "isolate":
		return okx.MgnModeIsolated, nil
	default:
		return "", fmt.Errorf("invalid mgnMode %q: want cross or isolated", s)
	}
}

// silentLogger discards all SDK logging. The SDK would otherwise echo request
// errors (which the cmd layer already surfaces cleanly) and info chatter to
// stdout, polluting output — especially under --json.
type silentLogger struct{}

func (silentLogger) Infof(string, ...any)  {}
func (silentLogger) Warnf(string, ...any)  {}
func (silentLogger) Debugf(string, ...any) {}
func (silentLogger) Errorf(string, ...any) {}
