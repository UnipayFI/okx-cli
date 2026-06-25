# OKX CLI

A command-line tool for the OKX API developed in Go. OKX runs a single
**unified trading account** (the v5 REST API, `/api/v5/*`), so one client serves
spot and every futures line; the product is chosen per command via the
instrument id (`instId`), instrument type (`instType`) and trade mode (`tdMode`).

Commands are grouped by product:

- **account** — identity & permissions, aggregate balance/equity, and account
  health (margin ratio).
- **spot** — per-coin assets and spot order entry (create / cancel / get / open).
- **futures** — account health, futures order entry and positions (list / close).

Authentication uses OKX's HMAC-SHA256 signing scheme (`OK-ACCESS-KEY` /
`OK-ACCESS-SIGN` / `OK-ACCESS-TIMESTAMP` / `OK-ACCESS-PASSPHRASE`). Only
authenticated (private) endpoints are covered — public market data is out of
scope.

## Installation and Configuration

### Install (prebuilt binary)
```shell
curl -sSL https://raw.githubusercontent.com/UnipayFI/okx-cli/refs/heads/main/download.sh | bash
```
Downloads the latest release for your platform/arch from GitHub Releases.

### Build from source
```shell
go build -o okx-cli .
```
Releases are produced by the `Release` GitHub Action (`.github/workflows/release.yml`),
which cross-compiles for Linux/macOS/Windows (amd64 + arm64) on every `v*` tag and
injects the version via ldflags.

### Environment variables
Before using, set your OKX API credentials (from the OKX API-management page):
```shell
export OKX_API_KEY="..."         # API key
export OKX_API_SECRET="..."      # API secret
export OKX_PASSPHRASE="..."      # passphrase set when the key was created

# Optional
export OKX_PROXY="socks5://127.0.0.1:1080"   # route REST traffic through a proxy
export OKX_BASE_URL="https://www.okx.com"    # override REST base URL
export OKX_DEMO="true"                        # use demo (paper) trading
```

> OKX rejects requests whose timestamp is skewed more than 30s from server
> time; the CLI best-effort syncs the server clock on every invocation.

### Output format
Every command supports a global `--json` flag. Without it, results render as a
table; with it, the raw API response is printed as indented JSON, e.g.:
```shell
./okx-cli account balance --json
```

## Usage

```
./okx-cli [command] [subcommand] [flags]
```

### account
```shell
./okx-cli account info                  # identity, account level, position mode, KYC, permissions
./okx-cli account balance               # aggregate equity/margin + per-coin balances (non-zero)
./okx-cli account health                # total/adjusted equity, UPL, IMR/MMR, margin ratio
```

### spot
```shell
./okx-cli spot assets                   # per-coin spot balances (non-zero)

# place a limit buy
./okx-cli spot order create -s BTC-USDT -S buy -t limit -q 0.001 -p 20000
# place a market buy spending 10 USDT
./okx-cli spot order create -s BTC-USDT -S buy -t market -q 10 -g quote_ccy

./okx-cli spot order open               # list open spot orders
./okx-cli spot order get    -s BTC-USDT -i <ordId>
./okx-cli spot order cancel -s BTC-USDT -i <ordId>
```

### futures
The product line is selected with the persistent `--instType` flag (default
`swap`); perpetual instruments are quoted `BASE-QUOTE-SWAP`, e.g. `BTC-USDT-SWAP`.
```shell
./okx-cli futures health                # account health (margin ratio & risk)

# open a cross-margin long (1 contract) at a limit price
./okx-cli futures order create -s BTC-USDT-SWAP -S buy -t limit -q 1 -p 20000 -m cross

./okx-cli futures order open                       # list open futures orders
./okx-cli futures order cancel -s BTC-USDT-SWAP -i <ordId>

./okx-cli futures position list                    # list open positions
./okx-cli futures position close -s BTC-USDT-SWAP -m cross   # market-close
```

## License
MIT — see [LICENSE](LICENSE).
