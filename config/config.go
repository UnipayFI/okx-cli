package config

import (
	"os"
	"strconv"
)

// Config holds the runtime configuration assembled from environment variables
// and global CLI flags. Credentials come from the OKX API-management page; all
// three (key, secret, passphrase) are required to sign private requests.
var Config struct {
	APIKey     string
	APISecret  string
	Passphrase string
	Proxy      string
	BaseURL    string
	Demo       bool
	OutputJSON bool
}

func init() {
	Config.APIKey = os.Getenv("OKX_API_KEY")
	Config.APISecret = os.Getenv("OKX_API_SECRET")
	Config.Passphrase = os.Getenv("OKX_PASSPHRASE")
	Config.Proxy = proxyFromEnv()
	Config.BaseURL = os.Getenv("OKX_BASE_URL")
	if v := os.Getenv("OKX_DEMO"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			Config.Demo = b
		}
	}
}

// proxyFromEnv resolves the proxy URL from the conventional proxy environment
// variables, with the same scheme-aware precedence curl/wget and Go's own
// http.ProxyFromEnvironment apply. OKX REST is HTTPS-only, so the scheme-specific
// HTTPS_PROXY wins, then the universal ALL_PROXY; for each variable the uppercase
// form is preferred over the lowercase. HTTP_PROXY is intentionally omitted: by
// convention it governs plain http:// traffic only, so it must not route the
// HTTPS OKX endpoint. Returns "" when none is set. Supported schemes: http,
// https, socks5, socks5h.
func proxyFromEnv() string {
	for _, key := range []string{
		"HTTPS_PROXY", "https_proxy",
		"ALL_PROXY", "all_proxy",
	} {
		if v := os.Getenv(key); v != "" {
			return v
		}
	}
	return ""
}
